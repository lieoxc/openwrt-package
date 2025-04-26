#include "serial.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <termios.h>
#include <unistd.h>
#include <errno.h>
#include <sys/select.h>
#include <sys/time.h>

// 全局串口文件描述符
static int serial_fd = -1;

// 串口初始化
int serial_init(const char *device, int baud_rate) {
    serial_fd = open(device, O_RDWR | O_NOCTTY | O_NDELAY);
    if (serial_fd == -1) {
        perror("无法打开串口设备");
        return -1;
    }
    
    struct termios options;
    tcgetattr(serial_fd, &options);
    
    // 设置波特率
    speed_t baud;
    switch (baud_rate) {
        case 9600:   baud = B9600;   break;
        case 19200:  baud = B19200;  break;
        case 38400:  baud = B38400;  break;
        case 57600:  baud = B57600;  break;
        case 115200: baud = B115200; break;
        default:     baud = B115200; break;
    }
    cfsetispeed(&options, baud);
    cfsetospeed(&options, baud);
    
    // 设置其他参数
    options.c_cflag |= (CLOCAL | CREAD);
    options.c_cflag &= ~PARENB;  // 无校验
    options.c_cflag &= ~CSTOPB;  // 1位停止位
    options.c_cflag &= ~CSIZE;   // 清除数据位设置
    options.c_cflag |= CS8;      // 8位数据位
    options.c_lflag &= ~(ICANON | ECHO | ECHOE | ISIG);  // 原始模式
    options.c_iflag &= ~(INPCK | ICRNL | IGNCR);         // 禁用输入处理
    options.c_oflag &= ~OPOST;                           // 禁用输出处理
    
    tcsetattr(serial_fd, TCSANOW, &options);
    return 0;
}

// 串口关闭
void serial_close(void) {
    if (serial_fd != -1) {
        close(serial_fd);
        serial_fd = -1;
    }
}

// 将十六进制字符串转换为字节数组
static int hex_to_bytes(const char *hex_string, uint8_t *bytes, int max_bytes) {
    int len = strlen(hex_string);
    int byte_count = 0;
    
    // 每两个字符转换为一个字节
    for (int i = 0; i < len && byte_count < max_bytes; i += 2) {
        if (i + 1 < len) {  // 确保有两个字符可以转换
            char hex_byte[3] = {hex_string[i], hex_string[i + 1], '\0'};
            bytes[byte_count++] = (uint8_t)strtol(hex_byte, NULL, 16);
        }
    }
    
    return byte_count;
}

// 串口发送十六进制字符串
int serial_send_hex(const char *hex_string) {
    if (serial_fd == -1) {
        printf("串口未初始化\n");
        return -1;
    }
    
    uint8_t bytes[256];
    int byte_count = hex_to_bytes(hex_string, bytes, sizeof(bytes));
    
    if (byte_count > 0) {
        // 打印发送的数据（调试用）
        printf("发送数据: ");
        for (int i = 0; i < byte_count; i++) {
            printf("0x%02x ", bytes[i]);
        }
        printf("\n");
        
        ssize_t written = write(serial_fd, bytes, byte_count);
        if (written < 0) {
            perror("串口发送失败");
            return -1;
        }
        printf("成功发送 %zd 字节数据\n", written);
        return written;
    }
    
    return -1;
}

// 串口接收单个字节（带超时）
int serial_receive_byte(char *byte, int timeout_ms) {
    if (serial_fd == -1) {
        printf("串口未初始化\n");
        return -1;
    }
    
    fd_set rfds;
    struct timeval tv;
    FD_ZERO(&rfds);
    FD_SET(serial_fd, &rfds);
    
    tv.tv_sec = timeout_ms / 1000;
    tv.tv_usec = (timeout_ms % 1000) * 1000;
    
    int retval = select(serial_fd + 1, &rfds, NULL, NULL, &tv);
    if (retval == -1) {
        perror("select 出错");
        return -1;
    } else if (retval) {
        if (FD_ISSET(serial_fd, &rfds)) {
            ssize_t bytes_read = read(serial_fd, byte, 1);
            if (bytes_read == 1) {
                return 1;
            }
        }
    }
    
    return 0;
}

// 串口接收Modbus-RTU数据包
int serial_receive_modbus(char *buffer, int buffer_size, int timeout_ms) {
    if (serial_fd == -1) {
        printf("串口未初始化\n");
        return -1;
    }
    
    int index = 0;
    struct timeval last_byte_time;
    gettimeofday(&last_byte_time, NULL);
    
    while (index < buffer_size) {
        char byte;
        int bytes_read = serial_receive_byte(&byte, timeout_ms);
        
        if (bytes_read == 1) {
            buffer[index++] = byte;
            gettimeofday(&last_byte_time, NULL);
            
            // 如果已经接收到至少4个字节（地址+功能码+CRC），检查是否接收完成
            if (index >= 4) {
                // 计算从最后一个字节到现在的时间（毫秒）
                struct timeval current_time;
                gettimeofday(&current_time, NULL);
                long time_diff = (current_time.tv_sec - last_byte_time.tv_sec) * 1000 +
                               (current_time.tv_usec - last_byte_time.tv_usec) / 1000;
                
                // 如果超过3.5个字符时间（在115200波特率下约为0.3ms），认为数据包接收完成
                if (time_diff > 1) {  // 使用1ms作为阈值
                    return index;
                }
            }
        } else if (bytes_read == 0) {
            // 超时，如果已经接收到数据，返回已接收的数据长度
            if (index > 0) {
                return index;
            }
            return 0;
        } else {
            return -1;  // 发生错误
        }
    }
    
    return index;  // 缓冲区已满
} 