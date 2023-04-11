package entity

import (
	"database/sql"
	"github.com/go-redis/redis"
)

type Connection struct {
	SqlDB   *sql.DB
	NoSqlDB *redis.Client
}
