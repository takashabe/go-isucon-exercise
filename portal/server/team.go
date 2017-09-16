package server

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/takashabe/go-isucon-exercise/portal/models"
)

const (
	sessionKeyTeam = "team_id"
)

func (s *Server) currentTeam(w http.ResponseWriter, r *http.Request) (*models.Team, error) {
	sess, err := s.session.SessionStart(w, r)
	if err != nil {
		return nil, err
	}

	item := sess.Get(sessionKeyTeam)
	id, ok := item.(int)
	if !ok {
		return nil, errors.New("invalid session id")
	}

	team, err := models.NewTeam().Get(id)
	if err != nil {
		sess.Delete(id)
		s.session.SessionDestroy(w, r)
		return nil, err
	}
	return team, nil
}

// unauthorized send response that Unauthorized status
func (s *Server) unauthorized(w http.ResponseWriter, err error) {
	Error(w, http.StatusUnauthorized, err, "failed to authentication")
}

// redirect redirect to location
func (s *Server) redirect(w http.ResponseWriter, r *http.Request, location string) {
	http.Redirect(w, r, location, http.StatusFound)
}

// Logout returns login page with cleanup session
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	s.session.SessionDestroy(w, r)
	JSON(w, http.StatusOK, "")
}

// Login authentication user and save user for session
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.unauthorized(w, err)
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	team, err := models.NewTeam().Authentication(email, password)
	if err != nil {
		s.unauthorized(w, err)
		return
	}

	sess, err := s.session.SessionStart(w, r)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, "failed to construct session")
	}
	sess.Set(sessionKeyTeam, team.ID)
	s.redirect(w, r, "/")
}
