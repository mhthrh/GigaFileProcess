package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mhthrh/GigaFileProcess/Api/View"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	ip          string
	port        int
	osInterrupt chan os.Signal
)

func init() {
	osInterrupt = make(chan os.Signal)
}
func main() {
	flag.StringVar(&ip, "ip listener", "localhost", "Mandalorian Episode")
	flag.IntVar(&port, "port listener", 8569, "Mandalorian Episode")

	flag.Parse()
	serverSync := http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: View.RunSync(),
	}

	if err := serverSync.ListenAndServe(); err != nil {
		log.Fatalf("cannot start service,%v ", err)
	}
	go signal.Notify(osInterrupt, os.Interrupt)
	<-osInterrupt

	if err := serverSync.Shutdown(context.Background()); err != nil {
		fmt.Printf("HTTP serverSync Shutdown: %v\n", err)
	}
}
