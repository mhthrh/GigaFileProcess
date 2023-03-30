package Test

import (
	"context"
	"github.com/mhthrh/GigaFileProcess/db"
	"github.com/mhthrh/GigaFileProcess/entity"
	"testing"
)

func TestDb(t *testing.T) {
	var trans []entity.FileStructure
	trans = append(trans, entity.FileStructure{
		ID:              0,
		FullName:        "",
		SourceIBAN:      "",
		DestinationIBAN: "",
		Amount:          0,
	})
	d := db.New(cnn)
	d.Insert(context.Background(), trans)
}
