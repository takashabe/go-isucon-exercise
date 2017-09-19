package server

import (
	"net/http"

	"github.com/takashabe/go-isucon-exercise/portal/models"
)

// History returns score histories
func (s *Server) History(w http.ResponseWriter, r *http.Request) {
	team, err := s.currentTeam(w, r)
	if err != nil {
		s.unauthorized(w, err)
		return
	}

	scores, err := models.NewScore().History(team.ID)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get history")
	}
	JSON(w, http.StatusOK, scores)
}

// ScoreDetail returns score detail
func (s *Server) ScoreDetail(w http.ResponseWriter, r *http.Request, id int) {
	team, err := s.currentTeam(w, r)
	if err != nil {
		s.unauthorized(w, err)
		return
	}

	score, err := models.NewScore().Get(id, team.ID)
	if err != nil {
		Error(w, http.StatusNotFound, err, "failed to get score")
	}
	JSON(w, http.StatusOK, score)
}
