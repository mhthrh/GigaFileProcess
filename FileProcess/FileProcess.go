package FileProcess

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mhthrh/GigaFileProcess/Validation"
	"github.com/mhthrh/GigaFileProcess/entity"
	"github.com/mhthrh/GigaFileProcess/rabbit"
	"strings"
)

const packages = 100_000

var (
	cancels  []context.CancelFunc
	channels []chan []string
	client   *redis.Client
	rabbit   *Rabbit.Mq
	i        = 1
	j        = packages
)

func New() error {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "", 0),
		Password: "r.RedisConnection.Password",
		DB:       0,
	})
	_, err := c.Ping().Result()
	if err != nil {
		return fmt.Errorf("cannot connet to redis,%w", err)
	}

	mq, err := Rabbit.New("")
	if err != nil {
		_ = c.Close()
		return fmt.Errorf("canot connet to rabbitMq server,%v", err)
	}
	client = c
	rabbit = mq
	return nil
}

func DoProcess(lines []string) error {
	var damegedList []string
	for {
		damagedLines := make(chan []string)
		ctx, can := context.WithCancel(context.Background())
		cancels = append(cancels, can)

		switch length := len(lines[i-1 : j]); {
		case length <= j:
			go process(lines, ctx, &damagedLines)
			channels = append(channels, damagedLines)
			goto exitFor
		default:
			go process(lines[i-1:j], ctx, &damagedLines)
			channels = append(channels, damagedLines)
			i = j
			j += j
		}
	}
exitFor:

	for _, channel := range channels {
		for _, s := range <-channel {
			damegedList = append(damegedList, s)
		}
	}
	return nil
}

func process(lines []string, ctx context.Context, dam *chan []string) {
	var damaged []string
	var obj []entity.FileStructure
	line := make(chan string)
	finish := make(chan struct{})
	_ = rabbit.DeclareQueue(rabbit.ID.String())
	go func() {
		for _, l := range lines {
			line <- l
		}
		finish <- struct{}{}
	}()

	for {
		select {
		case l := <-line:
			values := strings.Split(l, ",")
			if len(values) != 6 {
				damaged = append(damaged, fmt.Sprintf("%s#%s", l, "Array count is mismatch"))
				continue
			}
			newId, err := Validation.ValidaID(values[0])
			if err != nil {
				damaged = append(damaged, fmt.Sprintf("%s#%s", l, err.Error()))
				continue
			}
			err = Validation.ValidateFullName(values[1])
			if err != nil {
				damaged = append(damaged, fmt.Sprintf("%s#%s", l, err.Error()))
				continue
			}
			err = Validation.ValidateIban(values[2])
			if err != nil {
				damaged = append(damaged, fmt.Sprintf("%s#%s", l, err.Error()))
				continue
			}
			err = Validation.ValidateFullName(values[3])
			if err != nil {
				damaged = append(damaged, fmt.Sprintf("%s#%s", l, err.Error()))
				continue
			}
			amount, err := Validation.ValidateAmount(values[4])
			if err != nil {
				damaged = append(damaged, fmt.Sprintf("%s#%s", l, err.Error()))
				continue
			}
			obj = append(obj, entity.FileStructure{
				ID:              newId,
				FullName:        values[1],
				SourceIBAN:      values[2],
				DestinationIBAN: values[3],
				Amount:          amount,
			})
		case <-ctx.Done():
			return
		case <-finish:
			*dam <- damaged
			for _, o := range obj {
				byt, _ := json.Marshal(o)
				rabbit.Produce(rabbit.ID.String(), string(byt))
			}
		}
	}

}
