package rgmq

import (
	"errors"
	"github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal"
	"github.com/jackylee92/rgo/core/rgjson"
	"github.com/jackylee92/rgo/core/rglog"
	"github.com/jackylee92/rgo/core/rgrequest"
	"github.com/rs/xid"
	"github.com/streadway/amqp"
	"time"
)

// TODO <LiJunDong : 2022-06-04 00:36:32> --- 将vhost exchange queue 改为入参 支持监听发送多中队列
/*
 * @Content : rgmq
 * @Author  : LiJunDong
 * @Time    : 2022-05-27$
 */

const (
	configRabbitMQHeartBeat string = "util_rabbitmq_heartbeat" // rabbitmq 心跳
	configRabbitMQLog       string = "util_rabbitmq_log"       // rabbitmq 请求日志
	configRabbitMQChanMax   string = "util_rabbitmq_chan_max"  // rabbitmq 最大channel并发数
)

var chanMax = getChanMax()
var chanMaxLock = make(chan struct{}, chanMax)

type Client struct {
	this   *rgrequest.Client `json:"-"` // json解析忽略
	conn   *amqp.Connection  `json:"-"`
	config Config
}

type Config struct {
	Host     string
	Exchange string
	Routing  string
	Vhost    string
	Queue    string
	AutoAck  bool
}

// New 链接rabbitmq
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-05-27
func New(this *rgrequest.Client, config Config) (c *Client, err error) {
	conn, err := getConn(config)
	if err != nil {
		return c, err
	}
	c = new(Client)
	c.this = this
	c.config = config
	c.conn = conn
	return c, err
}

// getConn
// @Param   : config Config
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-09
func getConn(config Config) (conn *amqp.Connection, err error) {
	if config.Host == "" {
		return conn, errors.New("host不能为空")
	}
	if config.Vhost == "" {
		return conn, errors.New("vhost不能为空")
	}
	heartbeat := rgconfig.GetInt(configRabbitMQHeartBeat)
	if heartbeat == 0 {
		heartbeat = 30
	}
	conn, err = amqp.DialConfig(
		config.Host,
		amqp.Config{
			Vhost:      config.Vhost,
			Heartbeat:  time.Duration(heartbeat) * time.Second,
			ChannelMax: int(chanMax) + 50, // 允许的最大channel数量，比并发控制多50缓冲
		},
	)
	if err != nil {
		return conn, errors.New("rabbitmq链接失败|" + err.Error())
	}
	return conn, err
}

// Publish 发送消息
// @Param   : data 发送的数据， complete: 发送后mq回调结果执行函数
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-06-18
func (c *Client) Publish(data string, complete func(bool, string)) (err error) {
	if err = c.checkClient(data); err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			c.this.Log.Error("rabbitMq 捕获panic", err)
		}
	}()
	chanMaxLock <- struct{}{}
	defer func() {
		<-chanMaxLock
	}()
	messageId := xid.New().String()
	if rgconfig.GetBool(configRabbitMQLog) {
		logMap := map[string]interface{}{
			"vhost":    c.config.Vhost,
			"exchange": c.config.Exchange,
			"routing":  c.config.Routing,
			"queue":    c.config.Queue,
			"message":  messageId,
			"body":     data,
		}
		c.this.Log.Info("rabbitMq producer ", logMap)
	}
	if c.conn == nil || c.conn.IsClosed() {
		conn, err := getConn(c.config)
		if err != nil {
			return errors.New("rabbitMq 链接失败")
		}
		c.conn = conn
	}
	if c.conn == nil {
		return errors.New("rabbitmq connection链接为空")
	}
	channel, err := c.conn.Channel()
	if err != nil {
		return errors.New("rabbitmq channel链接失败|" + err.Error())
	}
	defer channel.Close()
	channel.Confirm(false)
	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1)) // 处理确认逻辑
	defer c.publishComplete(confirms, messageId, data, complete)       // 处理方法
	if err := channel.ExchangeDeclare(c.config.Exchange, amqp.ExchangeDirect, true, false, false, false, nil); err != nil {
		return errors.New("Exchange声明失败|" + err.Error())
	}
	if _, err := channel.QueueDeclare(c.config.Queue, true, false, false, false, nil); err != nil {
		return errors.New("Queue声明失败|" + err.Error())
	}
	if err := channel.QueueBind(c.config.Queue, c.config.Routing, c.config.Exchange, false, nil); err != nil {
		return errors.New("queue绑定exchange失败|" + err.Error())
	}
	err = channel.Publish(c.config.Exchange, c.config.Routing, false, false, amqp.Publishing{
		MessageId:   messageId,
		ContentType: "text/plain",
		Body:        []byte(data),
	})
	if err != nil {
		return errors.New("发送MQ消息失败|" + err.Error())
	}
	return err
}

// publishComplete 确认MQ消息发送成功
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-17
func (c *Client) publishComplete(confirms <-chan amqp.Confirmation, messageId string, data string, complete func(bool, string)) {
	// 消息确认
	result := false
	if confirmed := <-confirms; confirmed.Ack {
		result = true
		c.this.Log.Info("rabbitMq producer 消息确认成功", messageId)
	} else {
		c.this.Log.Info("rabbitMq producer 消息确认失败", messageId)
	}
	if complete != nil {
		complete(result, data)
	}
}

const (
	defaultChanMax = 10000
)

// getChanMax 获取配置的最大并发数
// @Param   :
// @Return  : data int
// @Author  : LiJunDong
// @Time    : 2022-05-28
func getChanMax() int {
	data := rgconfig.GetInt(configRabbitMQChanMax)
	if data == 0 {
		return defaultChanMax
	}
	return int(data)
}

// checkClient 检查push参数
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-05-28
func (c *Client) checkClient(data string) (err error) {
	if c.this == nil {
		return errors.New("this对象不能为空")
	}
	if data == "" {
		return errors.New("data消息不能为空")
	}
	return err
}

// Listen 开始监听
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-03
func (c *Client) Listen(pf func(delivery amqp.Delivery) bool) (err error) {
	if c.config.Queue == "" {
		return
	}
	if c.conn == nil || c.conn.IsClosed() {
		return errors.New("消费者创建MQ链接失败")
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return errors.New("消费者创建channel失败|" + err.Error())
	}
	if _, err := ch.QueueDeclare(c.config.Queue, true, false, false, false, nil); err != nil {
		return errors.New("消费者创建Queue声明失败|" + err.Error())
	}
	err = ch.Qos(2, 0, false)
	if err != nil {
		rglog.SystemError("消费者设置预加载条数失败" + err.Error())
		return
	}
	autoAck := c.config.AutoAck
	consumeId := rgglobal.AppName
	msgs, err := ch.Consume(
		c.config.Queue, // queue
		consumeId,      // consumer
		autoAck,        // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		rglog.SystemError("消费者监听队列失败|" + err.Error())
		return
	}
	logOpen := false
	if rgconfig.GetBool(configRabbitMQLog) {
		logOpen = true
	}
	// <LiJunDong : 2022-06-09 18:20:29> --- 阻塞
	for msg := range msgs {
		if logOpen {
			logData := map[string]interface{}{
				"host":         c.config.Host,
				"vhost":        c.config.Vhost,
				"exchange":     c.config.Exchange,
				"routing":      c.config.Routing,
				"queue":        c.config.Queue,
				"body":         string(msg.Body),
				"app_id":       msg.AppId,
				"consumer_tag": msg.ConsumerTag,
				"content_type": msg.ContentType,
				"message_id":   msg.MessageId,
				"type":         msg.Type,
				"user_id":      msg.UserId,
			}
			logDataJson, _ := rgjson.Marshel(logData)
			rgo.This.Log.Info("rabbitMq consumer", logDataJson)
		}
		ok := pf(msg)
		if ok && !autoAck {
			if logOpen {
				rgo.This.Log.Info("rabbitMq consumer ack true", msg.MessageId)
			}
			msg.Ack(ok)
		}
	}
	return err
}
