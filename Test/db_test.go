package Test

import (
	"context"
	"github.com/mhthrh/GigaFileProcess/Entity"
	"github.com/mhthrh/GigaFileProcess/db"
	"testing"
)

func TestDb(t *testing.T) {
	var trans []Entity.FileStructure
	trans = append(trans, Entity.FileStructure{
		ID:              0,
		FullName:        "",
		SourceIBAN:      "",
		DestinationIBAN: "",
		Amount:          0,
	})
	d := db.New(cnn)
	d.Insert(context.Background(), trans)
}
