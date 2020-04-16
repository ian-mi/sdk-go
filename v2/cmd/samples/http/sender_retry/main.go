package main

import (
	"context"
	"log"
	"time"

	cloudevents "github.com/ian-mi/sdk-go/v2"
)

func main() {
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://localhost:8080/")

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	// must send each event within 5 seconds for sleepy demo.

	log.Println("--- Constant ---")
	send10(cloudevents.ContextWithRetriesConstantBackoff(ctx, 10*time.Millisecond, 10), c)
	log.Println("--- Linear ---")
	send10(cloudevents.ContextWithRetriesLinearBackoff(ctx, 10*time.Millisecond, 10), c)
	log.Println("--- Exponential ---")
	send10(cloudevents.ContextWithRetriesExponentialBackoff(ctx, 10*time.Millisecond, 10), c)
}

func send10(ctx context.Context, c cloudevents.Client) {
	for i := 0; i < 100; i++ {
		e := cloudevents.NewEvent()
		e.SetType("com.cloudevents.sample.sent")
		e.SetSource("https://github.com/ian-mi/sdk-go/v2/cmd/samples/httpb/sender")
		_ = e.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
			"id":      i,
			"message": "Hello, World!",
		})

		if result := c.Send(ctx, e); !cloudevents.IsACK(result) {
			log.Printf("failed to send: %s", result.Error())
		} else {
			log.Printf("sent: %d", i)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
