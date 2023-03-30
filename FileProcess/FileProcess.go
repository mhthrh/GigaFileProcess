package FileProcess

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mhthrh/GigaFileProcess/Validation"
	"github.com/mhthrh/GigaFileProcess/entity"
	"github.com/mhthrh/GigaFileProcess/rabbit"
	"strings"
	"time"
)

const packages = 100_000

type Process struct {
	client  *redis.Client
	rabbit  *Rabbit.Mq
	invalid chan string
}

func New() (*Process, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "localhost", 6379),
		Password: "",
		DB:       0,
	})
	_, err := c.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("cannot connet to redis,%w", err)
	}

	mq, err := Rabbit.New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("canot connet to rabbitMq server,%v", err)
	}
	return &Process{
		client:  c,
		rabbit:  mq,
		invalid: make(chan string),
	}, nil
}

func (p *Process) DoProcess(lines []string) {

	_ = p.rabbit.DeclareQueue("rabbit.ID.String()")

	if len(lines) < packages {
		go process(lines, &p.invalid)
		goto wait
	}
	go process(lines[0:packages], &p.invalid)
	p.DoProcess(lines[packages:])
wait:
	var invalidList []string
	for line := range p.invalid {
		invalidList = append(invalidList, line)
	}
}

func process(lines []string, invalid *chan string) {
	line := make(chan string)
	finish := make(chan struct{})

	go func() {
		for _, l := range lines {
			line <- l
		}
		<-time.After(time.Millisecond * 200)
		finish <- struct{}{}
	}()

	for {
		select {
		case l := <-line:
			values := strings.Split(l, ",")
			if len(values) != 6 {
				*invalid <- fmt.Sprintf("%s#%s", l, "Array count is mismatch")
				continue
			}
			newId, err := Validation.ValidaID(values[0])
			if err != nil {
				*invalid <- fmt.Sprintf("%s#%s", l, err.Error())
				continue
			}
			err = Validation.ValidateFullName(values[1])
			if err != nil {
				*invalid <- fmt.Sprintf("%s#%s", l, err.Error())
				continue
			}
			err = Validation.ValidateIban(values[2])
			if err != nil {
				*invalid <- fmt.Sprintf("%s#%s", l, err.Error())
				continue
			}
			err = Validation.ValidateFullName(values[3])
			if err != nil {
				*invalid <- fmt.Sprintf("%s#%s", l, err.Error())
				continue
			}
			amount, err := Validation.ValidateAmount(values[4])
			if err != nil {
				*invalid <- fmt.Sprintf("%s#%s", l, err.Error())
				continue
			}

			byt, _ := json.Marshal(entity.FileStructure{
				ID:              newId,
				FullName:        values[1],
				SourceIBAN:      values[2],
				DestinationIBAN: values[3],
				Amount:          amount,
			})
			_ = rabbit.Produce("rabbit.ID.String()", string(byt))
		case <-finish:
			return
		}
	}

}
