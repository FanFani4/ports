package api

import (
	"github.com/FanFani4/ports/ports"
	"github.com/valyala/fasthttp"
)

const defaultLimit = 50

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"-"`
	Reason  string      `json:"reason,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func getResponse(data interface{}, success bool, code int, reason string) *Response {
	return &Response{
		Success: success,
		Code:    code,
		Reason:  reason,
		Data:    data,
	}
}

func (a *API) get(ctx *fasthttp.RequestCtx) *Response {
	portID := string(ctx.QueryArgs().Peek("id"))
	if portID == "" {
		return getResponse(nil, false, fasthttp.StatusBadRequest, "id is required")
	}

	port, err := a.cli.Get(ctx, &ports.GetArgs{Id: portID})
	if err != nil {
		return getResponse(nil, false, fasthttp.StatusInternalServerError, err.Error())
	}

	return getResponse(port, true, fasthttp.StatusOK, "")
}

func (a *API) list(ctx *fasthttp.RequestCtx) *Response {
	skip := int64(ctx.QueryArgs().GetUintOrZero("skip"))
	limit := int64(ctx.QueryArgs().GetUintOrZero("limit"))

	if limit == 0 {
		limit = defaultLimit
	}

	port, err := a.cli.List(ctx, &ports.ListArgs{Skip: skip, Limit: limit})
	if err != nil {
		return getResponse(nil, false, fasthttp.StatusInternalServerError, err.Error())
	}

	return getResponse(port, true, fasthttp.StatusOK, "")
}
