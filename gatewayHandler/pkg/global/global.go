package global

import (
	"gorm.io/gorm"
)

var VERSION = "0.0.5"
var VERSION_NUMBER = 5
var DB *gorm.DB
var OtaAddress string

type EventData struct {
	Name    string
	Message string
}

// 事件通道
var EventChan chan EventData
