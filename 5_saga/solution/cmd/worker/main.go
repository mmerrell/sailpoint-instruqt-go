package main

import (
	"log"

	"fulfillment"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{HostPort: "127.0.0.1:7233"})
	if err != nil {
		log.Fatalln("Unable to create client:", err)
	}
	defer c.Close()

	w := worker.New(c, fulfillment.TaskQueue, worker.Options{})

	w.RegisterWorkflow(fulfillment.FulfillmentWorkflow)
	w.RegisterActivity(&fulfillment.FulfillmentActivities{})

	log.Println("Worker started on task queue:", fulfillment.TaskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Worker error:", err)
	}
}
