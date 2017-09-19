package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// Datastore represent MySQL adapter
type Datastore struct {
	Conn *sql.DB
}

// NewDatastore returns initialized Datastore
func NewDatastore() (*Datastore, error) {
	db, err := sql.Open("mysql", "portal@tcp(localhost:3306)/portal?parseTime=true")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Datastore{Conn: db}, nil
}

func (d *Datastore) query(q string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.Conn.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args...)
}

func (d *Datastore) queryRow(q string, args ...interface{}) (*sql.Row, error) {
	stmt, err := d.Conn.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}

func (d *Datastore) findTeamByEmailAndPassword(email, password string) (*sql.Row, error) {
	return d.queryRow("select id, name, instance from teams where email=? and password=?",
		email, password)
}

func (d *Datastore) findTeamByID(id int) (*sql.Row, error) {
	return d.queryRow("select id, name, instance from teams where id=?", id)
}

func (d *Datastore) findQueueByTeamID(teamID int) (*sql.Row, error) {
	return d.queryRow("select msg_id from queues where team_id=?", teamID)
}

func (d *Datastore) saveQueues(teamID int, msgID string, submittedAt time.Time) error {
	stmt, err := d.Conn.Prepare("insert into queues (team_id, msg_id, submitted_at) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(teamID, msgID, submittedAt)
	return err
}

func (d *Datastore) saveScore(q QueueResponse) error {
	stmt, err := d.Conn.Prepare("insert into scores (team_id, summary, score, submitted_at, json) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var raw bytes.Buffer
	err = json.NewEncoder(&raw).Encode(q.BenchmarkResult)
	if err != nil {
		return err
	}

	summary, score := calculateScore(q.BenchmarkResult)
	_, err = stmt.Exec(q.TeamID, summary, score, q.CreatedAt, raw.String())
	return err
}

func (d *Datastore) findScoreByIDAndTeamID(scoreID, teamID int) (*sql.Row, error) {
	return d.queryRow("select id, summary, score, json from scores where id=? and team_id=?",
		scoreID, teamID)
}

func (d *Datastore) findScoreHistoryByIDAndTeamID(teamID int) (*sql.Rows, error) {
	return d.query("select id, summary, score, json from scores where team_id=? order by submitted_at DESC", teamID)
}
