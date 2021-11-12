package eth

import (
	"fmt"
	"io"
)

type LogDecoder struct {
	logEvent     *Log
	topicDecoder *Decoder
	DataDecoder  *Decoder

	topicIndex int
}

func NewLogDecoder(logEvent *Log) *LogDecoder {
	decoder := &LogDecoder{
		logEvent:     logEvent,
		topicDecoder: NewDecoder(nil),
	}

	if len(logEvent.Data) > 0 {
		decoder.DataDecoder = NewDecoder(logEvent.Data)
	}

	return decoder
}

func (d *LogDecoder) ReadTopic() ([]byte, error) {
	if d.topicIndex >= len(d.logEvent.Topics) {
		return nil, io.EOF
	}

	topic := d.logEvent.Topics[d.topicIndex]
	d.topicIndex++

	return topic, nil
}

func (d *LogDecoder) ReadTypedTopic(typeName string) (out interface{}, err error) {
	topic, err := d.ReadTopic()
	if err != nil {
		return nil, fmt.Errorf("read topic: %w", err)
	}

	return d.topicDecoder.SetBytes(topic).Read(typeName)
}

func (d *LogDecoder) ReadData(typeName string) (out interface{}, err error) {
	return d.DataDecoder.Read(typeName)
}
