package main

type WorkOrder struct {
	runningTime int
	tasks       []Task
}

// Master is manages the benchmark
type Master struct {
	ctx Ctx
}

// create Master and benchmark context
func NewMaster(host string, port int, file string, agent string) (*Master, error) {
	ctx := newCtx()
	if host != "" {
		ctx.host = host
	}
	if port != 0 {
		ctx.port = port
	}
	if file != "" {
		ctx.paramFile = file
	}
	if agent != "" {
		ctx.agent = agent
	}
	err := ctx.setupSessions()
	if err != nil {
		return nil, err
	}

	return &Master{ctx: *ctx}, nil
}

func (m *Master) start() ([]byte, error) {
	// TODO
	// 1. create workers
	// 2. run for each workers with order()
	// 3. sum return results from worker.run

	result := newResult()
	orders := IsuconWorkOrder()
	for _, o := range orders {
		w := NewWorker(m.ctx, o.runningTime, o.tasks)
		result.Merge(*w.run())
		if !result.Valid {
			break
		}
	}

	json, err := result.json()
	if err != nil {
		PrintDebugf("failed to result.json(): %s", err.Error())
		return nil, err
	}
	return json, nil
}
