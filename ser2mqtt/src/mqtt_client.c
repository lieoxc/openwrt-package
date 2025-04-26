#include "mqtt_client.h"
#include "serial.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "cJSON.h"

// 全局 MQTT 客户端实例
static struct mosquitto *mosq = NULL;

// MQTT连接回调
static void on_connect(struct mosquitto *mosq, void *userdata, int result) {
    if (result == 0) {
        printf("成功连接到MQTT服务器\n");
        // 订阅 /data/recv
        mosquitto_subscribe(mosq, NULL, "/data/recv", 0);
    } else {
        fprintf(stderr, "MQTT连接失败\n");
        exit(EXIT_FAILURE);
    }
}

// MQTT消息回调
static void on_message(struct mosquitto *mosq, void *userdata, const struct mosquitto_message *msg) {
    printf("收到MQTT消息，主题: %s, 内容: %s\n", msg->topic, (char *)msg->payload);
    
    // 检查消息主题
    if (strcmp(msg->topic, "/data/recv") == 0) {
        // 解析JSON消息
        cJSON *json = cJSON_Parse((char *)msg->payload);
        if (json) {
            cJSON *msg_obj = cJSON_GetObjectItem(json, "msg");
            if (msg_obj && msg_obj->valuestring) {
                // 将消息发送到串口
                if (serial_send_hex(msg_obj->valuestring) < 0) {
                    printf("串口发送失败\n");
                }
            }
            cJSON_Delete(json);
        }
    }
}

// MQTT客户端初始化
void mqtt_init(void) {
    const char *host = "192.168.6.175";  // MQTT服务器地址
    const int port = 1883;
    const char *username = "admin";  // 用户名
    const char *password = "qazwsx@1234";  // 密码
    
    mosquitto_lib_init();

    mosq = mosquitto_new(NULL, true, NULL);
    if (!mosq) {
        fprintf(stderr, "创建MQTT实例失败\n");
        exit(EXIT_FAILURE);
    }

    mosquitto_username_pw_set(mosq, username, password);
    mosquitto_connect_callback_set(mosq, on_connect);
    mosquitto_message_callback_set(mosq, on_message);

    if (mosquitto_connect(mosq, host, port, 60) != MOSQ_ERR_SUCCESS) {
        fprintf(stderr, "连接MQTT服务器失败\n");
        mosquitto_destroy(mosq);
        mosquitto_lib_cleanup();
        exit(EXIT_FAILURE);
    }
    
    mosquitto_loop_start(mosq);
}

// MQTT客户端清理
void mqtt_deinit(void) {
    if (mosq) {
        mosquitto_destroy(mosq);
    }
    mosquitto_lib_cleanup();
}

// MQTT发布消息
void mqtt_publish(const char *topic, const char *message) {
    if (mosquitto_publish(mosq, NULL, topic, strlen(message), message, 0, false) != MOSQ_ERR_SUCCESS) {
        fprintf(stderr, "MQTT发布失败 [%s]: %s\n", topic, message);
    } else {
        printf("MQTT发布成功 [%s]: %s\n", topic, message);
    }
}

// 获取MQTT客户端实例
struct mosquitto *get_mqtt_client(void) {
    return mosq;
} 