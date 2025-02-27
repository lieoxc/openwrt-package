package privateMqtt

import (
	"encoding/json"
	"fmt"
	"gatewayhandler/model"
	config "gatewayhandler/mqtt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var privateMqttClient mqtt.Client

func CreateMqttClient() {
	// 初始化配置
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.PrivateMqttConfig.Broker)
	opts.SetUsername(config.PrivateMqttConfig.User)
	opts.SetPassword(config.PrivateMqttConfig.Pass)
	opts.SetClientID("weather-Station")
	// 干净会话
	opts.SetCleanSession(true)
	// 恢复客户端订阅，需要broker支持
	opts.SetResumeSubs(true)
	// 自动重连
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(20 * time.Second)
	// 消息顺序
	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logrus.Debug("mqtt connect success")
	})
	// 断线重连
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logrus.Error("mqtt connect  lost: ", err)
		privateMqttClient.Disconnect(250)
		// 等待连接成功，失败重新连接
		for {
			token := privateMqttClient.Connect()
			if token.Wait() && token.Error() == nil {
				fmt.Println("Reconnected to MQTT broker")
				break
			}
			fmt.Printf("Reconnect failed: %v\n", token.Error())
			time.Sleep(5 * time.Second)
		}
	})

	privateMqttClient = mqtt.NewClient(opts)
	for {
		if token := privateMqttClient.Connect(); token.Wait() && token.Error() != nil {
			logrus.Error("MQTT Broker 1 连接失败:", token.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}
func makeGatewayPubPayload(cfgID, devID string, payload []byte) ([]byte, error) {
	var dataMap = make(map[string]interface{})
	err := json.Unmarshal(payload, &dataMap)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	gatewayData := make(map[string]interface{})
	subDeviceData := make(map[string]map[string]interface{})

	if cfgID == model.DefaultGatewayCfgID { // 网关自身
		gatewayData = dataMap
	} else { // 子设备
		subDeviceData[devID] = dataMap
	}
	// 创建结构体实例
	publish := model.GatewayPublish{
		GatewayData:   &gatewayData,
		SubDeviceData: &subDeviceData,
	}
	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(publish, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func makeGatewayEventPayload(cfgID, devID string, payload []byte) ([]byte, error) {
	var dataMap model.EventInfo
	err := json.Unmarshal(payload, &dataMap)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	var gatewayData model.EventInfo
	subDeviceData := make(map[string]model.EventInfo)

	if cfgID == model.DefaultGatewayCfgID { // 网关自身
		gatewayData = dataMap
	} else { // 子设备
		subDeviceData[devID] = dataMap
	}
	// 创建结构体实例
	publish := model.GatewayCommandPulish{
		GatewayData:   &gatewayData,
		SubDeviceData: &subDeviceData,
	}
	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(publish, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// 上报telemetry消息
func ForwardTelemetryMessage(cfgID, devID string, payload []byte) error {
	qos := byte(config.PrivateMqttConfig.Telemetry.QoS)
	pubTelemetryTopic := config.PrivateMqttConfig.Telemetry.GatewayPublishTopic + "/" + cfgID + "/" + devID
	// 转发消息
	jsonData, err := makeGatewayPubPayload(cfgID, devID, payload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("privateMqttClient pub topic:%v payload%v", pubTelemetryTopic, string(jsonData))
	token := privateMqttClient.Publish(pubTelemetryTopic, qos, false, jsonData)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

// 上报attributes消息
func ForwardAttributesMessage(cfgID, devID string, payload []byte) error {
	qos := byte(config.PrivateMqttConfig.Telemetry.QoS)
	pubAttributesTopic := config.PrivateMqttConfig.Attributes.GatewayPublishTopic + "/" + cfgID + "/" + devID
	logrus.Info("topic:", pubAttributesTopic, "value:", string(payload))
	// 发布消息
	jsonData, err := makeGatewayPubPayload(cfgID, devID, payload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("privateMqttClient pub topic:%v payload%v", pubAttributesTopic, string(jsonData))

	token := privateMqttClient.Publish(pubAttributesTopic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

// 上报tevents消息
func ForwardEventsMessage(cfgID, devID string, payload []byte) error {
	qos := byte(config.PrivateMqttConfig.Telemetry.QoS)
	pubEventsTopic := config.PrivateMqttConfig.Events.GatewayPublishTopic + "/" + cfgID + "/" + devID
	logrus.Info("topic:", pubEventsTopic, "value:", string(payload))
	// 发布消息
	jsonData, err := makeGatewayEventPayload(cfgID, devID, payload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("privateMqttClient pub topic:%v payload%v", pubEventsTopic, string(jsonData))
	token := privateMqttClient.Publish(pubEventsTopic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}
