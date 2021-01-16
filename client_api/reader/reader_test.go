package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/FanFani4/ports/ports"
	"github.com/bcicen/jstream"
	"github.com/sirupsen/logrus"
)

var testBuf = []byte(`{
  "AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "Abu Dhabi",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  }
}`)

var invalidBuf = []byte(`"AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  }`)

func TestJSONReader_GetPorts(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *JSONReader
		inspect func(r *JSONReader, t *testing.T) //inspects receiver after test run

		want string
	}{
		{
			name: "ok",
			want: `[{"id":"AEAJM","name":"Ajman","city":"Ajman","country":"United Arab Emirates","coordinates":[55.513645,25.405216],"province":"Ajman","timezone":"Asia/Dubai","unlocs":["AEAJM"],"code":"52000"},{"id":"AEAUH","name":"Abu Dhabi","city":"Abu Dhabi","country":"United Arab Emirates","coordinates":[54.37,24.47],"province":"Abu Z¸aby [Abu Dhabi]","timezone":"Asia/Dubai","unlocs":["AEAUH"],"code":"52001"}]`,
			init: func(t *testing.T) *JSONReader {
				buf := bytes.NewBuffer(testBuf)
				return &JSONReader{
					ctx:     context.TODO(),
					ports:   make(chan *ports.Port),
					decoder: jstream.NewDecoder(buf, 1).EmitKV(),
					log:     logrus.New(),
				}
			},
		},
		{
			name: "fail",
			want: `[]`,
			init: func(t *testing.T) *JSONReader {
				buf := bytes.NewBuffer(invalidBuf)
				return &JSONReader{
					ctx:     context.TODO(),
					ports:   make(chan *ports.Port),
					decoder: jstream.NewDecoder(buf, 1).EmitKV(),
					log:     logrus.New(),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			got1 := receiver.GetPorts()

			gotSlice := []*ports.Port{}

			for item := range got1 {
				gotSlice = append(gotSlice, item)
			}

			b, _ := json.Marshal(gotSlice)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(string(b), tt.want) {
				t.Errorf("JSONReader.GetPorts got1 = %v, want1: %v", string(b), tt.want)
			}
		})
	}
}
