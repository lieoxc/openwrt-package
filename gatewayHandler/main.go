package main

import (
	"flag"
	"gatewayhandler/initialize"
	mqttapp "gatewayhandler/mqtt"
	"gatewayhandler/mqtt/localMqtt"
	"gatewayhandler/mqtt/privateMqtt"
	"gatewayhandler/register"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func main() {
	// 1. 定义命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/conf.yml", "Path to config file")
	flag.Parse()

	initialize.ViperInit(configPath)
	initialize.LogInIt()
	_, err := initialize.PgInit(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	err = mqttapp.MqttInit()
	if err != nil {
		logrus.Fatal(err)
	}

	// 创建一个MQTT客户端，并且订阅本地遥测，时间，属性转发消息
	localMqtt.SubscribeInit()
	if err != nil {
		logrus.Fatal(err)
	}
	// 创建一个MQTT客户端，连接到内网监控站MQTT服务器
	privateMqtt.CreateMqttClient()

	//	初始化设备注册
	register.DevRegisterInit()

	gracefulShutdown()
}

func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Println("gatewayHandler exiting")
}
