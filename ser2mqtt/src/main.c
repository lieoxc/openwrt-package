#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <net/if.h>
#include <sys/ioctl.h>
#include <arpa/inet.h>

#include "modbus.h"
#include "mqtt_client.h"
#include "serial.h"

// 获取MAC地址
char *get_mac_address(const char *interface) {
    static char mac_str[18];  // 用于存储 MAC 地址字符串
    struct ifreq ifr;
    int sock = socket(AF_INET, SOCK_DGRAM, 0);
    if (sock == -1) {
        perror("无法创建套接字");
        return NULL;
    }

    strncpy(ifr.ifr_name, interface, IFNAMSIZ - 1);
    if (ioctl(sock, SIOCGIFHWADDR, &ifr) == -1) {
        perror("无法获取 MAC 地址");
        close(sock);
        return NULL;
    }

    unsigned char *mac = (unsigned char *)ifr.ifr_hwaddr.sa_data;
    snprintf(mac_str, sizeof(mac_str), "%02x:%02x:%02x:%02x:%02x:%02x",
             mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);

    close(sock);
    return mac_str;
}

int main(int argc, char *argv[]) {
    // 初始化 MQTT
    mqtt_init();
    
    // 获取 eth0 的 MAC 地址
    char *mac_address = get_mac_address("eth0");
    if (!mac_address) {
        fprintf(stderr, "无法获取 MAC 地址\n");
        return -1;
    }
    
    // 初始化串口
    if (serial_init("/dev/ttyS0", 115200) != 0) {
        fprintf(stderr, "串口初始化失败\n");
        return -1;
    }
    
    char buf[255];

    while (1) {
        // 接收Modbus-RTU数据包
        int bytes_read = serial_receive_modbus(buf, sizeof(buf), 1000);  // 1秒超时
        if (bytes_read > 0) {
            // 检查Modbus-RTU格式
            if (check_modbus_rtu_format((uint8_t *)buf, bytes_read) == 0) {
                // 将接收到的字节数据转换为十六进制字符串
                char hex_string[512];
                for (int i = 0; i < bytes_read; i++) {
                    snprintf(hex_string + i * 2, sizeof(hex_string) - i * 2, "%02x", (unsigned char)buf[i]);
                }
                
                printf("收到有效的Modbus-RTU数据: %s\n", hex_string);
                // 构造 JSON 数据
                char message[512];
                snprintf(message, sizeof(message), "{\"msg\":\"%s\"}", hex_string);
                char topic[64];
                snprintf(topic, sizeof(topic), "/ser2mqtt/report/%s", mac_address);
                // 将数据通过 MQTT 发布
                mqtt_publish(topic, message);
            } else {
                printf("收到无效的Modbus-RTU数据\n");
            }
        }
    }

    serial_close();
    mqtt_deinit();
    return 0;
} 