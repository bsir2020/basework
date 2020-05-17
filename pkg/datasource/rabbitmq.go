package datasource

import (
	"encoding/json"
	"fmt"
	"github.com/bsir2020/basework/api"
	cfg "github.com/bsir2020/basework/configs"
	"github.com/bsir2020/basework/pkg/log"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
)

// 定义生产者接口
type Producer interface {
	MsgContent(dataByte []byte) string
}

// 定义接收者接口
type Receiver interface {
	Consumer([]byte) error
}

// 定义RabbitMQ对象
type RabbitMQ struct {
	Url       string
	Qgame     string
	Qgamefeed string
	//user         string
	//password     string
	//ip           string //mq ip
	//port         int    //port
	//vhost        string // vhost
	connection *amqp.Connection
	channel    *amqp.Channel
	//queueName    string // 队列名称
	//routingKey   string // key名称
	//exchangeName string // 交换机名称
	//exchangeType string // 交换机类型
	producerList []Producer
	receiverList []Receiver
	mu           sync.RWMutex
}

// 定义队列交换机对象
type Queue struct {
	Url       string
	Qgame     string
	Qgamefeed string
	//User     string // mq user
	//Password string // mq password
	//Ip       string // mq ip
	//Port     int    // mq port
	//Vhost    string // vhost
	//QuName   string // 队列名称
	//RtKey    string // key值
	//ExName   string // 交换机名称
	//ExType   string // 交换机类型
}

// 链接rabbitMQ
func (r *RabbitMQ) mqConnect() {
	//RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.user, r.password, r.ip, r.port, r.vhost)
	var err error
	r.connection, err = amqp.Dial(r.Url)
	if err != nil {
		panic(err.Error())
	}

	r.channel, err = r.connection.Channel() // 赋值给RabbitMQ对象
	if err != nil {
		panic(err.Error())
	}
}

// 关闭RabbitMQ连接
func (r *RabbitMQ) mqClose() {
	// 先关闭管道,再关闭链接
	err := r.channel.Close()
	if err != nil {
		fmt.Printf("MQ管道关闭失败:%s \n", err)
	}
	err = r.connection.Close()
	if err != nil {
		fmt.Printf("MQ链接关闭失败:%s \n", err)
	}
}

var logger *zap.Logger

// 创建一个新的操作对象
func New() *RabbitMQ {
	logger = log.New()
	return &RabbitMQ{
		Url:       cfg.EnvConfig.Rabbitmq.Url,
		Qgame:     cfg.EnvConfig.Rabbitmq.Qgame,
		Qgamefeed: cfg.EnvConfig.Rabbitmq.Qgamefeed,
		//user:         q.User,
		//password:     q.Password,
		//port:         q.Port,
		//ip:           q.Ip,
		//vhost:        q.Vhost,
		//queueName:    q.QuName,
		//routingKey:   q.RtKey,
		//exchangeName: q.ExName,
		//exchangeType: q.ExType,
	}
}

// 启动RabbitMQ客户端,并初始化
func (r *RabbitMQ) Start() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)

			if r.connection != nil {
				r.connection.Close()
			}

			if r.channel != nil {
				r.channel.Close()
			}
		}
	}()

	r.initMQ()

	// 开启监听生产者发送任务
	//for _, producer := range r.producerList {
	//	go r.listenProducer(producer)
	//}
	// 开启监听接收者接收任务
	for _, receiver := range r.receiverList {
		go r.listenReceiver(receiver)
	}
}

func (r *RabbitMQ) initMQ() {
	// 验证链接是否正常,否则重新链接
	if r.channel == nil {
		r.mqConnect()
	}

	/*
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		err := r.channel.ExchangeDeclare(r.exchangeName, r.exchangeType, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册交换机失败:%s \n", err)
			return
		}

		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.channel.QueueDeclare(r.queueName, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return
		}

		// 队列绑定
		err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true, nil)
		if err != nil {
			fmt.Printf("MQ绑定队列失败:%s \n", err)
			return
		}
	*/
}

// 注册发送指定队列指定路由的生产者
func (r *RabbitMQ) RegisterProducer(producer Producer) {
	r.producerList = append(r.producerList, producer)
}

// 发送任务
//func (r *RabbitMQ) listenProducer(producer Producer) {
//	// 发送任务消息
//	err := r.channel.Publish(r.exchangeName, r.routingKey, false, false, amqp.Publishing{
//		ContentType: "application/json",
//		Body:        []byte(producer.MsgContent()),
//	})
//	if err != nil {
//		fmt.Printf("MQ任务发送失败:%s \n", err)
//		return
//	}
//}

func (r *RabbitMQ) SendGameFeed(msg []byte) {
	// 发送任务消息
	err := r.channel.Publish(r.Qgamefeed, r.Qgamefeed, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(msg),
	})
	if err != nil {
		fmt.Printf("MQ任务发送失败:%s \n", err)
		return
	}
}

// 注册接收指定队列指定路由的数据接收者
func (r *RabbitMQ) RegisterReceiver(receiver Receiver) {
	r.mu.Lock()
	r.receiverList = append(r.receiverList, receiver)
	r.mu.Unlock()
}

// 监听接收者接收任务
func (r *RabbitMQ) listenReceiver(receiver Receiver) {
	// 处理结束关闭链接
	//defer r.mqClose()

	// 获取消费通道,确保rabbitMQ一个一个发送消息
	//err := r.channel.Qos(1, 0, true)
	msgList, err := r.channel.Consume(r.Qgame, "", false, false, false, false, nil)
	if err != nil {
		//fmt.Printf("获取消费通道异常:%s \n", err)
		logger.Error(err.Error())
		return
	}

	for msg := range msgList {
		rpmsg := &api.Message{}
		json.Unmarshal(msg.Body, rpmsg)

		// 处理数据
		err := receiver.Consumer(msg.Body)
		if err != nil {
			//fmt.Printf("确认消息未完成异常:%s \n", err)
			logger.Error(err.Error())
			//r.channel.Nack(msg.DeliveryTag, true, true)
			err = msg.Ack(false)
			if err != nil {
				//fmt.Printf("确认消息完成异常:%s \n", err)
				logger.Error(err.Error())

				rpmsg.Status = false
			}
		} else {
			rpmsg.Status = true

			// 确认消息,必须为false
			err = msg.Ack(false)
			if err != nil {
				//fmt.Printf("确认消息完成异常:%s \n", err)
				logger.Error(err.Error())
			}
		}

		//回复
		if d, err := json.Marshal(rpmsg); err == nil {
			r.SendGameFeed(d)
		} else {
			logger.Error(err.Error(), zap.String("reply", string(msg.Body)))
		}
	}
}
