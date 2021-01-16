package api

import (
	"encoding/json"

	"github.com/FanFani4/ports/ports"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func NewAPI(log *logrus.Logger, cli ports.PortDomainServiceClient) *API {
	return &API{
		log: log,
		cli: cli,
	}
}

type API struct {
	log *logrus.Logger
	cli ports.PortDomainServiceClient
}

func (a *API) writeResponse(ctx *fasthttp.RequestCtx, response *Response) {
	ctx.SetStatusCode(response.Code)
	ctx.SetContentType("application/json")

	body, err := json.Marshal(response)
	if err != nil {
		a.log.Error("failed to marshal response: ", err)
		return
	}

	ctx.SetBody(body)
}

func (a *API) HandleFasthttp(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/list":
		a.writeResponse(ctx, a.list(ctx))
	case "/get":
		a.writeResponse(ctx, a.get(ctx))
	default:
		a.writeResponse(ctx, getResponse(nil, false, fasthttp.StatusNotFound, "not found"))
	}
}
