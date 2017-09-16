package server

import (
	"context"
	"net/http"

	"github.com/takashabe/go-isucon-exercise/portal/models"
)

// Queues returns active queue list
func (s *Server) Queues(w http.ResponseWriter, r *http.Request) {
	_, err := s.currentTeam(w, r)
	if err != nil {
		s.unauthorized(w, err)
		return
	}

	// TODO: implements model method
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
	err = q.Publish(context.Background(), team.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, "failed to enqueue")
	}

	JSON(w, http.StatusOK, "")
}
