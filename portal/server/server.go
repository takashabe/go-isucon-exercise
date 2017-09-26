package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/takashabe/go-router"
	session "github.com/takashabe/go-session"
	_ "github.com/takashabe/go-session/memory" // session driver
)

// ErrorResponse is Error response template
type ErrorResponse struct {
	Message string `json:"reason"`
	Error   error  `json:"-"`
}

func (e *ErrorResponse) String() string {
	return fmt.Sprintf("reason: %s, error: %v", e.Message, e.Error)
}

// Respond is response write to ResponseWriter
func Respond(w http.ResponseWriter, code int, src interface{}) {
	var body []byte
	var err error

	switch s := src.(type) {
	case []byte:
		if !json.Valid(s) {
			Error(w, http.StatusInternalServerError, err, "invalid json")
			return
		}
		body = s
	case string:
		body = []byte(s)
	case *ErrorResponse, ErrorResponse:
		// avoid infinite loop
		if body, err = json.Marshal(src); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"reason\":\"failed to parse json\"}"))
			return
		}
	default:
		if body, err = json.Marshal(src); err != nil {
			Error(w, http.StatusInternalServerError, err, "failed to parse json")
			return
		}
	}
	w.WriteHeader(code)
	w.Write(body)
}

// Error is wrapped Respond when error response
func Error(w http.ResponseWriter, code int, err error, msg string) {
	e := &ErrorResponse{
		Message: msg,
		Error:   err,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	Respond(w, code, e)
}

// JSON is wrapped Respond when success response
func JSON(w http.ResponseWriter, code int, src interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	Respond(w, code, src)
}

// Server supply HTTP server of the portal
type Server struct {
	port       int
	session    *session.Manager
	pubsubAddr string
}

// NewServer returns initialized Server
func NewServer(pubsubAddr string) (*Server, error) {
	session, err := session.NewManager("memory", "portal", 3600)
	if err != nil {
		return nil, err
	}
	return &Server{
		session:    session,
		pubsubAddr: pubsubAddr,
	}, nil
}

// Routes returns router
func (s *Server) Routes() *router.Router {
	r := router.NewRouter()

	// login page
	r.Get("/api/login", s.Logout)
	r.Post("/api/login", s.Login)

	// main page
	r.Get("/api/team", s.GetTeam)
	r.Get("/api/queues", s.Queues)
	r.Post("/api/enqueue", s.Enqueue)
	r.Get("/api/history", s.History)
	r.Get("/api/bench_detail/:id", s.ScoreDetail)
	// r.Get("/leader_board", nil)

	// frontend
	r.ServeFile("/", "./public/index.html")
	return r
}

// Run start server
func (s *Server) Run(port int) error {
	log.Println("starting server...")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.Routes())
}
