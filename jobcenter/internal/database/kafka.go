package database

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// KafkaConfig 定义Kafka客户端的配置参数
type KafkaConfig struct {
	Addr          string `json:"addr,optional"`          // Kafka服务器地址
	WriteCap      int    `json:"writeCap,optional"`      // 写入通道容量
	ReadCap       int    `json:"readCap,optional"`       // 读取通道容量
	ConsumerGroup string `json:"ConsumerGroup,optional"` // 消费者组ID
}

// KafkaData 定义Kafka消息的数据结构
type KafkaData struct {
	Topic string // 主题
	Key   []byte // 消息键
	Data  []byte // 消息内容
}

// KafkaClient Kafka客户端结构体
type KafkaClient struct {
	w         *kafka.Writer  // Kafka写入器
	r         *kafka.Reader  // Kafka读取器
	readChan  chan KafkaData // 读取通道
	writeChan chan KafkaData // 写入通道
	c         KafkaConfig    // 配置信息
	closed    bool           // 关闭标志
	mutex     sync.Mutex     // 互斥锁，用于保护关闭操作
}

// NewKafkaClient 创建新的Kafka客户端实例
func NewKafkaClient(c KafkaConfig) *KafkaClient {
	return &KafkaClient{
		c: c,
	}
}

// StartWrite 启动Kafka写入功能
// 初始化写入器和写入通道，并启动后台发送协程
func (k *KafkaClient) StartWrite() {
	w := &kafka.Writer{
		Addr:     kafka.TCP(k.c.Addr),
		Balancer: &kafka.LeastBytes{}, // 使用最少字节负载均衡策略
	}
	k.w = w
	k.writeChan = make(chan KafkaData, k.c.WriteCap)
	go k.sendKafka()
}

// Send 发送消息到Kafka
// 将消息放入写入通道，如果发生panic则标记客户端为已关闭
func (w *KafkaClient) Send(data KafkaData) {
	defer func() {
		if err := recover(); err != nil {
			w.closed = true
		}
	}()
	w.writeChan <- data
	w.closed = false
}

// Close 关闭Kafka客户端
// 关闭写入器和读取器，并清理相关资源
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

// sendKafka 后台协程，负责从写入通道读取消息并发送到Kafka
// 包含重试机制和错误处理
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
				// 尝试发送消息
				err = w.w.WriteMessages(ctx, messages...)
				if err == nil {
					success = true
					break
				}
				// 处理特定错误类型，进行重试
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
			// 如果发送失败，重新放入通道等待重试
			if !success {
				w.Send(data)
			}
		}
	}
}

// StartRead 启动Kafka读取功能
// 初始化读取器和读取通道，并启动后台读取协程
func (k *KafkaClient) StartRead() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.c.Addr},
		GroupID:  k.c.ConsumerGroup,
		MinBytes: 10e3, // 最小读取字节数：10KB
		MaxBytes: 10e6, // 最大读取字节数：10MB
	})
	k.r = r
	k.readChan = make(chan KafkaData, k.c.ReadCap)
	go k.readMsg()
}

// readMsg 后台协程，负责从Kafka读取消息并放入读取通道
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

// Read 从读取通道获取消息
// 返回一个KafkaData类型的消息
func (k *KafkaClient) Read() KafkaData {
	msg := <-k.readChan
	return msg
}
