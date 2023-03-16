package main

import "fmt"

func main() {
	fmt.Sprintf("amqp://%s:%s@%s:%d/", "cfg.Rabbit.UserName", "cfg.Rabbit.Password", "cfg.Rabbit.Host", "cfg.Rabbit.Port")
}
