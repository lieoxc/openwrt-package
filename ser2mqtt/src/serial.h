#ifndef SERIAL_H
#define SERIAL_H

#include <stdint.h>

// 串口初始化
int serial_init(const char *device, int baud_rate);

// 串口关闭
void serial_close(void);

// 串口发送十六进制字符串
int serial_send_hex(const char *hex_string);

// 串口接收数据（带超时）
int serial_receive(char *buffer, int buffer_size);

// 串口接收单个字节（带超时）
int serial_receive_byte(char *byte, int timeout_ms);

// 串口接收Modbus-RTU数据包
int serial_receive_modbus(char *buffer, int buffer_size, int timeout_ms);

#endif // SERIAL_H 