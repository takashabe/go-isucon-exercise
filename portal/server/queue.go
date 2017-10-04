package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/takashabe/go-isucon-exercise/portal/models"
)

// Queues returns active queue list
func (s *Server) Queues(w http.ResponseWriter, r *http.Request) {
	team, err := s.currentTeam(w, r)
	if err != nil {
		s.unauthorized(w, err)
		return
	}

	q, err := models.NewQueue(s.pubsubAddr)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, "failed to initialized queue")
		return
	}
	queues, err := q.CurrentQueues(context.Background(), team.ID)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get current queues status")
		return
	}

	JSON(w, http.StatusOK, queues)
}

// Enqueue send queue to the pubsub server
func (s *Server) Enqueue(w http.ResponseWriter, r *http.Request) {
	team, err := s.currentTeam(w, r)
	if err != nil {
		s.unauthorized(w, err)
		return
	}

	q, err := models.NewQueue(s.pubsubAddr)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, "fialed to initialized queue")
		return
	}
	_, err = q.Publish(context.Background(), team.ID)
	if err != nil {
		if err == models.ErrExistQueue {
			Error(w, http.StatusNotFound, err, err.Error())
			return
		}
		Error(w, http.StatusInternalServerError, err, "failed to enqueue")
		return
	}

	JSON(w, http.StatusOK, "")
}

var runPolling = false

// Polling call pull API from the result pubsub
func (s *Server) Polling() error {
	if runPolling {
		return nil
	}
	q, err := models.NewQueue(s.pubsubAddr)
	if err != nil {
		return err
	}

	runPolling = true
	go func() {
		for {
			ctx := context.Background()
			err = q.PullAndSave(ctx)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	return nil
}
