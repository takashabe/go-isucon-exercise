package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// Datastore represent MySQL adapter
type Datastore struct {
	conn *sql.DB
}

func newDatastore() (*Datastore, error) {
	db, err := sql.Open("mysql", "portal@tcp(localhost:3306)/portal?parseTime=true")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Datastore{conn: db}, nil
}

func (d *Datastore) query(q string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.conn.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args)
}

func (d *Datastore) queryRow(q string, args ...interface{}) (*sql.Row, error) {
	stmt, err := d.conn.Prepare(q)
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
