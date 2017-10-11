package benchmark

import (
	"net/http/httptest"
	"strconv"
	"strings"
)

type Helper struct{}

var helper = Helper{}

func (h *Helper) setAddr(ts *httptest.Server, ctx *Ctx) *Ctx {
	addr := strings.Split(ts.Listener.Addr().String(), ":")
	ctx.host = addr[0]
	ctx.port, _ = strconv.Atoi(addr[1])
	return ctx
}

func (h *Helper) testCtx(ts *httptest.Server) Ctx {
	ctx := newCtx()
	ctx.setupSessions()
	h.setAddr(ts, ctx)
	return *ctx
}

func (h *Helper) testDriver(ctx Ctx) *Driver {
	return &Driver{
		result: newResult(),
		ctx:    ctx,
	}
}
