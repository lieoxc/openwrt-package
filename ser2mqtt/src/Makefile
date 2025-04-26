LDFLAGS = -lmosquitto
TARGET = ser2mqtt

# 源文件和目标文件
SRCS = main.c modbus.c cJSON.c mqtt_client.c serial.c
OBJS = $(SRCS:.c=.o)

# 默认目标
all: $(TARGET)

# 生成可执行文件
$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^ $(LDFLAGS)

# 生成目标文件
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# 清理编译结果
clean:
	rm -f $(TARGET) $(OBJS)
