package database

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"sync"
	"time"
)

type KafkaConfig struct {
	Addr          string `json:"addr,optional"`
	WriteCap      int    `json:"writeCap,optional"`
	ReadCap       int    `json:"readCap,optional"`
	ConsumerGroup string `json:"ConsumerGroup,optional"`
}
type KafkaData struct {
	Topic string
	Key   []byte
	Data  []byte
}
type KafkaClient struct {
	w         *kafka.Writer
	r         *kafka.Reader
	readChan  chan KafkaData
	writeChan chan KafkaData
	c         KafkaConfig
	closed    bool
	mutex     sync.Mutex
}

func NewKafkaClient(c KafkaConfig) *KafkaClient {
	return &KafkaClient{
		c: c,
	}
}

func (k *KafkaClient) StartWrite() {
	w := &kafka.Writer{
		Addr:     kafka.TCP(k.c.Addr),
		Balancer: &kafka.LeastBytes{},
	}
	k.w = w
	k.writeChan = make(chan KafkaData, k.c.WriteCap)
	go k.sendKafka()
}

func (w *KafkaClient) Send(data KafkaData) {
	defer func() {
		if err := recover(); err != nil {
			w.closed = true
		}
	}()
	w.writeChan <- data
	w.closed = false
}

func (w *KafkaClient) Close() {
	if w.w != nil {
		w.w.Close()
		w.mutex.Lock()
		defer w.mutex.Unlock()
		if !w.closed {
			close(w.writeChan)
			w.closed = true
		}
	}
	if w.r != nil {
		w.r.Close()
	}
}

func (w *KafkaClient) sendKafka() {
	for {
		select {
		case data := <-w.writeChan:
			messages := []kafka.Message{
				{
					Topic: data.Topic,
					Key:   data.Key,
					Value: data.Data,
				},
			}
			var err error
			const retries = 3
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			success := false
			for i := 0; i < retries; i++ {
				// attempt to create topic prior to publishing the message
				err = w.w.WriteMessages(ctx, messages...)
				if err == nil {
					success = true
					break
				}
				if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					success = false
					continue
				}
				if err != nil {
					success = false
					log.Printf("kafka send writemessage err %s \n", err.Error())
				}
			}
			if !success {
				//重新放进去等待消费
				w.Send(data)
			}
		}
	}

}

func (k *KafkaClient) StartRead(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.c.Addr},
		Topic:    topic,
		GroupID:  k.c.ConsumerGroup,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	k.r = r
	k.readChan = make(chan KafkaData, k.c.ReadCap)
	go k.readMsg()
}

func (k *KafkaClient) StartReadNew(topic string) *KafkaClient {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.c.Addr},
		Topic:    topic,
		GroupID:  k.c.ConsumerGroup,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	client := NewKafkaClient(k.c)
	client.r = r
	client.readChan = make(chan KafkaData, k.c.ReadCap)
	go client.readMsg()
	return client
}

func (k *KafkaClient) readMsg() {
	for {
		m, err := k.r.ReadMessage(context.Background())
		if err != nil {
			logx.Error(err)
			continue
		}
		data := KafkaData{
			Key:  m.Key,
			Data: m.Value,
		}
		k.readChan <- data
	}
}

func (k *KafkaClient) Read() KafkaData {
	msg := <-k.readChan
	return msg
}

func (k *KafkaClient) Rput(data KafkaData) {
	k.readChan <- data
}

func (k *KafkaClient) SendSync(data KafkaData) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP(k.c.Addr),
		Balancer: &kafka.LeastBytes{},
	}
	w.AllowAutoTopicCreation = true
	messages := []kafka.Message{
		{
			Topic: data.Topic,
			Key:   data.Key,
			Value: data.Data,
		},
	}
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = w.WriteMessages(ctx, messages...)
	return err
}
