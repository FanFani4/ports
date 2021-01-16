package sender

import (
	"context"

	"github.com/FanFani4/ports/ports"
	"github.com/sirupsen/logrus"
)

type Sender struct {
	cli ports.PortDomainServiceClient
	ctx context.Context
	log *logrus.Logger
}

func (s *Sender) Send(ports <-chan *ports.Port) {
	for port := range ports {
		s.log.Info(port)
		resp, err := s.cli.Insert(s.ctx, port)
		if err != nil {
			s.log.Error("failed to send port: ", port.Id, err)
			continue
		}

		if !resp.Success {
			s.log.Error("failed to insert port: ", port.Id)
		}
	}
}

func NewSender(ctx context.Context, log *logrus.Logger, cli ports.PortDomainServiceClient) *Sender {
	return &Sender{cli, ctx, log}
}
