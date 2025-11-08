package main

import (
	"log"
	"time"

	"github.com/bruli/pinger/pkg/events"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = nc.Drain()
		nc.Close()
	}()

	resource := "testing"
	subject := events.PingSubjet

	var n int
	i := 50
	for range i {
		now := timestamppb.New(time.Now().UTC())
		msg := events.PingResult{
			Resource:  resource,
			Status:    "Failed",
			Latency:   float32(100 + n),
			CreatedAt: now,
		}

		data, err := proto.Marshal(&msg)
		if err != nil {
			log.Fatal(err)
		}
		if err = nc.Publish(subject, data); err != nil {
			log.Fatal(err)
		}
		n++
		time.Sleep(time.Second)
	}

	var count int

	no := 20
	for range no {
		readyMsg := events.PingResult{
			Resource:  resource,
			Status:    "ok",
			Latency:   float32(250 + count),
			CreatedAt: timestamppb.New(time.Now().UTC()),
		}

		data, err := proto.Marshal(&readyMsg)
		if err != nil {
			log.Fatal(err)
		}

		if err = nc.Publish(subject, data); err != nil {
			log.Fatal(err)
		}
		count++
		time.Sleep(time.Second)
	}

	nc.Close()
}
