#include "modbus.h"
#include <stdio.h>

// CRC16 计算函数
uint16_t calculate_crc16(const uint8_t *data, size_t length) {
    uint16_t crc = 0xFFFF;
    
    for (size_t i = 0; i < length; i++) {
        crc ^= (uint16_t)data[i];
        for (int j = 0; j < 8; j++) {
            if (crc & 0x0001) {
                crc = (crc >> 1) ^ 0xA001;
            } else {
                crc = crc >> 1;
            }
        }
    }
    
    return crc;
}

// 检查Modbus-RTU数据格式
int check_modbus_rtu_format(const uint8_t *data, size_t length) {
    // 检查最小长度（地址 + 功能码 + CRC = 4字节）
    if (length < 4) {
        printf("数据长度不足\n");
        return -1;
    }
    
    // 检查功能码
    uint8_t function_code = data[1];
    if (function_code != MODBUS_READ_COILS &&
        function_code != MODBUS_READ_DISCRETE &&
        function_code != MODBUS_READ_HOLDING &&
        function_code != MODBUS_READ_INPUT &&
        function_code != MODBUS_WRITE_COIL &&
        function_code != MODBUS_WRITE_REGISTER &&
        function_code != MODBUS_WRITE_COILS &&
        function_code != MODBUS_WRITE_REGISTERS) {
        printf("无效的功能码: 0x%02X\n", function_code);
        return -1;
    }
    
    // 计算CRC
    uint16_t calculated_crc = calculate_crc16(data, length - 2);
    uint16_t received_crc = (data[length - 1] << 8) | data[length - 2];
    
    if (calculated_crc != received_crc) {
        printf("CRC校验失败: 计算值=0x%04X, 接收值=0x%04X\n", calculated_crc, received_crc);
        return -1;
    }
    
    return 0;
} 