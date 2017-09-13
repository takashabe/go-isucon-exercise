package portal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/takashabe/go-router"
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
	PrintDebugf("%v", e)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	Respond(w, code, e)
}

// JSON is wrapped Respond when success response
func JSON(w http.ResponseWriter, code int, src interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	Respond(w, code, src)
}

// Routes returns router
func Routes() *router.Router {
	r := router.NewRouter()

	// login page
	r.Get("/login", nil)
	r.Post("/login", nil)

	// main page
	r.Get("/", nil)
	r.Get("/messages", nil)
	r.Get("/queues", nil)
	r.Post("/enqueue", nil)
	r.Get("/history", nil)
	r.Get("/bench_detail:id", nil)
	r.Get("/leader_board", nil)

	return r
}

// Server supply HTTP server of the portal
type Server struct {
	// TODO: queue server configuration
}

// Run start server
func (s *Server) Run(port int) error {
	log.Println("starting server...")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), Routes())
}
