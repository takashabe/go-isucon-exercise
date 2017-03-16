package main

import "time"

type WorkOrder struct {
	runningTime time.Duration
	tasks       []Task
}

// Master is manages the benchmark
type Master struct {
	ctx Ctx
}

// create Master and benchmark context
func NewMaster(host string, port int, file string, agent string) (*Master, error) {
	ctx := newCtx()
	if host != defaultHost {
		ctx.host = host
	}
	if port != defaultPort {
		ctx.port = port
	}
	if file != defaultFile {
		ctx.paramFile = file
	}
	if agent != defaultAgent {
		ctx.agent = agent
	}
	err := ctx.setupSessions()
	if err != nil {
		return nil, err
	}

	return &Master{ctx: *ctx}, nil
}

func (m *Master) start() (string, error) {
	result := newResult()
	orders := IsuconWorkOrder()
	for _, o := range orders {
		w := NewWorker(m.ctx, o.runningTime, o.tasks)
		PrintDebugf("RUN worker: %v\n", w)
		result.Merge(*w.run())
		if !result.Valid {
			break
		}
	}

	json, err := result.json()
	if err != nil {
		PrintDebugf("failed to result.json(): %s", err.Error())
		return "", err
	}
	return string(json), nil
}
