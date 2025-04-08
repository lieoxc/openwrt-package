#!/bin/sh

# 日志文件路径
LOG_FILE="/var/log/plugin-monitor.log"

# 日志记录函数
log_message() {
    local message="$1"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $message" >> "$LOG_FILE"
}

# 检查单个插件的心跳
check_plugin_heartbeat() {
    local plugin_name=$1
    if [ "$plugin_name" = "iotService" ]; then
        local heartbeat_file="/tmp/iot_heartbeat"
    else
        local heartbeat_file="/tmp/${plugin_name}_heartbeat"
    fi
    local timeout=30  # 5分钟超时
    
    #log_message "正在检查插件 ${plugin_name} 的心跳状态..."
    
    if [ ! -f "$heartbeat_file" ]; then
        log_message "警告：插件 ${plugin_name} 的心跳文件不存在"
        log_message "插件 ${plugin_name} 心跳超时，正在重启..."
        /etc/init.d/${plugin_name} restart
        return
    fi
    
    local current_time=$(date +%s)
    local heartbeat_time=$(stat -c %Y "$heartbeat_file")
    local time_diff=$((current_time - heartbeat_time))
    
    #log_message "插件 ${plugin_name} 心跳信息："
    #log_message "  当前时间：$(date -d @$current_time '+%Y-%m-%d %H:%M:%S')"
    #log_message "  最后心跳：$(date -d @$heartbeat_time '+%Y-%m-%d %H:%M:%S')"
    #log_message "  时间差：${time_diff}秒"
    
    if [ $time_diff -gt $timeout ]; then
        log_message "警告：插件 ${plugin_name} 心跳超时（${time_diff}秒 > ${timeout}秒）"
        log_message "插件 ${plugin_name} 心跳超时，正在重启..."
        /etc/init.d/${plugin_name} restart
    fi
}

# 监控所有插件
monitor_plugins() {
    for plugin in iotService data_collect; do
        check_plugin_heartbeat "$plugin"
    done
}

# 主循环
while true; do
    monitor_plugins
    sleep 10
done 