package main

import (
	"flag"
	"gatewayhandler/initialize"
	mqttapp "gatewayhandler/mqtt"
	"gatewayhandler/mqtt/publish"
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

	err := mqttapp.MqttInit()
	if err != nil {
		logrus.Fatal(err)
	}
	publish.CreateMqttClient()
	if err != nil {
		logrus.Fatal(err)
	}

	register.DevRegisterInit()

	gracefulShutdown()
}

func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Println("gatewayHandler exiting")
}
