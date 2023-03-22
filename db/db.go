package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mhthrh/GigaFileProcess/Entity"
	"strings"
)

var (
	packages []string
)

const pack = 10_000

type Database struct {
	db *sql.DB
}

func New(d *sql.DB) *Database {
	return &Database{db: d}
}
func (d *Database) Insert(ctx context.Context, trans []Entity.FileStructure) error {
	makePackages(trans)
	if len(packages) < 1 {
		return fmt.Errorf("transaction slice is empty")
	}
	_, err := d.db.ExecContext(ctx, "insert into transactions (id,fullname,s_iban,d_iban,amount,datetime) %s")
	if err != nil {
		return fmt.Errorf("cannot insert to db, %v", err)
	}
	return nil
}

func makePackages(trans []Entity.FileStructure) {
	if len(trans) < pack {
		packages = append(packages, doProcess(trans))
		return
	}
	packages = append(packages, doProcess(trans[0:pack]))
	makePackages(trans[pack:])
}

func doProcess(trans []Entity.FileStructure) string {
	var sb strings.Builder
	for _, tran := range trans {
		sb.WriteString(fmt.Sprintf("select %d,%s,%s,%s,%2f from dual union all", tran.ID, tran.FullName, tran.SourceIBAN, tran.DestinationIBAN, tran.Amount))
	}
	return sb.String()
}
