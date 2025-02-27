package localMqtt

import (
	"strings"
	"time"

	config "gatewayhandler/mqtt"
	"gatewayhandler/mqtt/privateMqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

var LocalMqttClient mqtt.Client

func SubscribeInit() error {
	// 创建mqtt客户端
	subscribeMqttClient()
	//订阅attribute、event、telemetry、command、ota升级消息, 从路由器本地转发过来
	err := subscribe()
	return err
}

func subscribe() error {
	// 订阅attribute消息
	err := SubscribeAttribute()
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 订阅event消息
	err = SubscribeEvent()
	if err != nil {
		logrus.Error(err)
		return err
	}
	//订阅telemetry消息
	err = SubscribeTelemetry()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func subscribeMqttClient() {
	// 初始化配置
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.LocalMqttConfig.Broker)
	opts.SetUsername(config.LocalMqttConfig.User)
	opts.SetPassword(config.LocalMqttConfig.Pass)
	id := "localMqtt-sub-" + uuid.New()[0:8]
	opts.SetClientID(id)
	logrus.Info("clientid: ", id)

	// 干净会话
	opts.SetCleanSession(true)
	// 恢复客户端订阅，需要broker支持
	opts.SetResumeSubs(true)
	// 自动重连
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(200 * time.Second)
	// 消息顺序
	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logrus.Println("mqtt connect success")
	})
	// 断线重连
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logrus.Println("mqtt connect  lost: ", err)
		LocalMqttClient.Disconnect(250)
		for {
			if token := LocalMqttClient.Connect(); token.Wait() && token.Error() != nil {
				logrus.Error("MQTT Broker 1 连接失败:", token.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			subscribe()
			break
		}
	})

	LocalMqttClient = mqtt.NewClient(opts)
	// 等待连接成功，失败重新连接
	for {
		if token := LocalMqttClient.Connect(); token.Wait() && token.Error() != nil {
			logrus.Error("MQTT Broker 1 连接失败:", token.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

}

// 订阅telemetry消息
func SubscribeTelemetry() error {

	deviceTelemetryMessageHandler := func(_ mqtt.Client, d mqtt.Message) {
		datas := strings.Split(string(d.Topic()), "/")
		if len(datas) != 4 {
			logrus.Error("topic error:", string(d.Topic()))
			return
		}
		cfgID, devID := datas[2], datas[3]
		privateMqtt.ForwardTelemetryMessage(cfgID, devID, d.Payload())
	}

	topic := config.LocalMqttConfig.Telemetry.SubscribeTopic
	logrus.Info("Local subscribe topic:", topic)
	qos := byte(config.LocalMqttConfig.Telemetry.QoS)
	if token := LocalMqttClient.Subscribe(topic, qos, deviceTelemetryMessageHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

// 订阅attribute消息，暂不需要线程池，不需要消息队列
func SubscribeAttribute() error {
	// 订阅attribute消息
	deviceAttributeHandler := func(_ mqtt.Client, d mqtt.Message) {
		datas := strings.Split(string(d.Topic()), "/")
		if len(datas) != 4 {
			logrus.Error("topic error:", string(d.Topic()))
			return
		}
		cfgID, devID := datas[2], datas[3]
		privateMqtt.ForwardAttributesMessage(cfgID, devID, d.Payload())
	}
	topic := config.LocalMqttConfig.Attributes.SubscribeTopic
	logrus.Info("subscribe topic:", topic)
	qos := byte(config.LocalMqttConfig.Attributes.QoS)
	if token := LocalMqttClient.Subscribe(topic, qos, deviceAttributeHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

// 订阅event消息，暂不需要线程池，不需要消息队列
func SubscribeEvent() error {
	// 订阅event消息
	deviceEventHandler := func(_ mqtt.Client, d mqtt.Message) {
		datas := strings.Split(string(d.Topic()), "/")
		if len(datas) != 4 {
			logrus.Error("topic error:", string(d.Topic()))
			return
		}
		cfgID, devID := datas[2], datas[3]
		privateMqtt.ForwardEventsMessage(cfgID, devID, d.Payload())
	}
	topic := config.LocalMqttConfig.Events.SubscribeTopic
	qos := byte(config.LocalMqttConfig.Events.QoS)
	if token := LocalMqttClient.Subscribe(topic, qos, deviceEventHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}
