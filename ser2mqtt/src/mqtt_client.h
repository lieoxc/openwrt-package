#ifndef MQTT_CLIENT_H
#define MQTT_CLIENT_H

#include <mosquitto.h>

// MQTT客户端初始化
void mqtt_init(void);

// MQTT客户端清理
void mqtt_deinit(void);

// MQTT发布消息
void mqtt_publish(const char *topic, const char *message);

// 获取MQTT客户端实例
struct mosquitto *get_mqtt_client(void);

#endif // MQTT_CLIENT_H 