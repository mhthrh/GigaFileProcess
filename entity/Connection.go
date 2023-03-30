package entity

import (
	"database/sql"
	"github.com/go-redis/redis"
	Rabbit "github.com/mhthrh/GigaFileProcess/rabbit"
)

type Connection struct {
	SqlDB   *sql.DB
	NoSqlDB *redis.Client
	MQ      *Rabbit.Mq
}
