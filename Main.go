package main

import "fmt"

func main() {
	fmt.Sprintf("amqp://%s:%s@%s:%d/", "cfg.rabbit.UserName", "cfg.rabbit.Password", "cfg.rabbit.Host", "cfg.rabbit.Port")
}
