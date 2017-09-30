package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/pkg/errors"
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
	db.SetMaxIdleConns(32)
	db.SetMaxOpenConns(32)

	return &Datastore{Conn: db}, nil
}

// Close calls DB.Close
func (d *Datastore) Close() error {
	return d.Conn.Close()
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

func (d *Datastore) saveTeams(id int, name, email, password, instance string) error {
	stmt, err := d.Conn.Prepare("insert into teams (id, name, email, password, instance) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, email, password, instance)
	return err
}

func (d *Datastore) findQueueByTeamID(teamID int) (*sql.Row, error) {
	return d.queryRow("select msg_id, finished_at from queues where team_id=? order by id desc", teamID)
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

func (d *Datastore) saveScoreAndUpdateQueue(q QueueResponse) error {
	var raw bytes.Buffer
	err := json.NewEncoder(&raw).Encode(q.BenchmarkResult)
	if err != nil {
		return err
	}

	tx, err := d.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	update, err := tx.Prepare("update queues set finished_at=NOW() where msg_id=?")
	if err != nil {
		return err
	}
	defer update.Close()

	insert, err := tx.Prepare("insert into scores (team_id, summary, score, submitted_at, json) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer insert.Close()

	res, err := update.Exec(q.SourceMessageID)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("not exist updated rows")
	}

	summary, score := calculateScore(q.BenchmarkResult)
	res, err = insert.Exec(q.TeamID, summary, score, q.CreatedAt, raw.String())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Datastore) findScoreByIDAndTeamID(scoreID, teamID int) (*sql.Row, error) {
	return d.queryRow("select id, summary, score, json, submitted_at from scores where id=? and team_id=?", scoreID, teamID)
}

func (d *Datastore) findScoreHistoryByIDAndTeamID(teamID int) (*sql.Rows, error) {
	return d.query("select id, summary, score, submitted_at from scores where team_id=? order by submitted_at DESC", teamID)
}
