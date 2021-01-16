package reader

import (
	"context"
	"os"

	"github.com/FanFani4/ports/ports"
	"github.com/bcicen/jstream"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type JSONReader struct {
	ctx     context.Context
	ports   chan *ports.Port
	decoder *jstream.Decoder
	log     *logrus.Logger
}

func NewJSONReader(ctx context.Context, log *logrus.Logger, jsonPath string) (*JSONReader, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	reader := &JSONReader{
		ctx:     ctx,
		ports:   make(chan *ports.Port),
		decoder: jstream.NewDecoder(file, 1).EmitKV(),
		log:     log,
	}

	return reader, nil
}

func (j *JSONReader) GetPorts() <-chan *ports.Port {
	go j.readPorts()

	return j.ports
}

func (j *JSONReader) readPorts() {
	stream := j.decoder.Stream()

	for {
		select {
		case mv := <-stream:
			if mv == nil {
				close(j.ports)
				return
			}

			port := ports.Port{}

			val, ok := mv.Value.(jstream.KV)
			if !ok {
				continue
			}

			err := mapstructure.Decode(val.Value, &port)
			if err != nil {
				continue
			}

			port.Id = val.Key
			j.ports <- &port
		case <-j.ctx.Done():
			close(j.ports)
			return
		}
	}
}
