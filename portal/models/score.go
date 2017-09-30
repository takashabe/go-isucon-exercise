package models

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const tooSlowResponseMessage = "ミリ秒以内に応答しませんでした"

func calculateScore(raw BenchmarkResult) (summary string, score int) {
	baseScore := raw.Response.Success + int(float64(raw.Response.Redirect)*0.1)
	minusScore := (raw.Response.ServerError * 10) + (raw.Response.Exception * 20)
	tooSlowPenalty := 0
	for _, v := range raw.Violations {
		if strings.Contains(v.Cause, tooSlowResponseMessage) {
			tooSlowPenalty += v.Count
		}
	}
	tooSlowPenalty *= 100

	score = baseScore - minusScore - tooSlowPenalty
	if score < 0 {
		score = 0
	}
	if raw.Valid && score > 1 {
		summary = "success"
	} else {
		summary = "fail"
	}
	return
}

// Score represent benchmark score
type Score struct {
	ID          int       `json:"id"`
	Summary     string    `json:"summary"`
	Score       int       `json:"score"`
	Detail      string    `json:"detail"`
	SubmittedAt timeStamp `json:"submitted_at"`
}

type timeStamp struct {
	mysql.NullTime
}

func (t *timeStamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

// NewScore returns initialized Score
func NewScore() *Score {
	return &Score{}
}

// Get returns score object
func (s *Score) Get(scoreID, teamID int) (*Score, error) {
	d, err := NewDatastore()
	if err != nil {
		return nil, err
	}
	defer d.Close()

	row, err := d.findScoreByIDAndTeamID(scoreID, teamID)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&s.ID, &s.Summary, &s.Score, &s.Detail, &s.SubmittedAt)
	return s, err
}

// History returns all team score
func (s *Score) History(teamID int) ([]*Score, error) {
	d, err := NewDatastore()
	if err != nil {
		return nil, err
	}
	defer d.Close()

	rows, err := d.findScoreHistoryByIDAndTeamID(teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	history := []*Score{}
	for rows.Next() {
		score := &Score{}
		err := rows.Scan(&score.ID, &score.Summary, &score.Score, &score.SubmittedAt)
		if err != nil {
			return nil, err
		}
		history = append(history, score)
	}
	return history, nil
}
