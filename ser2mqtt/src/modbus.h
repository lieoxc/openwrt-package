#ifndef MODBUS_H
#define MODBUS_H

#include <stdint.h>
#include <stddef.h>

// Modbus-RTU 功能码定义
#define MODBUS_READ_COILS          0x01
#define MODBUS_READ_DISCRETE       0x02
#define MODBUS_READ_HOLDING        0x03
#define MODBUS_READ_INPUT          0x04
#define MODBUS_WRITE_COIL         0x05
#define MODBUS_WRITE_REGISTER     0x06
#define MODBUS_WRITE_COILS        0x0F
#define MODBUS_WRITE_REGISTERS    0x10

// CRC16 计算函数
uint16_t calculate_crc16(const uint8_t *data, size_t length);

// 检查Modbus-RTU数据格式
int check_modbus_rtu_format(const uint8_t *data, size_t length);

#endif // MODBUS_H 