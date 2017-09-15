package models

import "strings"

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
