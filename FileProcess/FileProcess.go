package FileProcess

import (
	"context"
	"fmt"
	"github.com/mhthrh/GigaFileProcess/Validation"
	"strings"
)

const packages = 100_000

var (
	cancels  []context.CancelFunc
	channels []chan []string
	i        = 1
	j        = packages
)

type MyStruct struct {
	ID              int64
	FullName        string
	SourceIBAN      string
	DestinationIBAN string
	Amount          float32
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
	var obj []MyStruct
	line := make(chan string)
	finish := make(chan struct{})
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
			obj = append(obj, MyStruct{
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
		}
	}

}
