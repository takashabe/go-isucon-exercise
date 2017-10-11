package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/takashabe/go-isucon-exercise/portal/models"
)

// Exec queue.Publish() by any teamID
func main() {
	var (
		teamID     int
		pubsubAddr string
	)
	flag.IntVar(&teamID, "team_id", 1, "team_id")
	flag.StringVar(&pubsubAddr, "pubsub", "http://localhost:9000", "pubsub server url")
	flag.Parse()

	q, err := models.NewQueue(pubsubAddr)
	if err != nil {
		panic(err)
	}

	msgID, err := q.Publish(context.Background(), teamID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("published message: %s", msgID)
}
