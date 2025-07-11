#!/bin/sh

# 删除telemetry_datas表中指定时间之前的数据的脚本

# 数据库连接参数 - 根据您的信息配置
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="thingspanel"
DB_USER="postgres"
DB_PASSWORD="postgresThingsPanel"

# 检查是否安装了psql
if ! command -v psql &> /dev/null
then
    echo "错误：psql 未安装，请先安装PostgreSQL客户端"
    echo "在Ubuntu/Debian上可以运行: sudo apt-get install postgresql-client"
    echo "在CentOS/RHEL上可以运行: sudo yum install postgresql"
    exit 1
fi

# 使用 TRUNCATE 清空表（速度快且立即释放空间）
START_TIME=$(date +%s)

PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
TRUNCATE TABLE public.telemetry_datas;
VACUUM FULL public.telemetry_datas;
TRUNCATE TABLE public.command_set_logs;
VACUUM FULL public.command_set_logs;
EOF

if [[ $? -ne 0 ]]; then
    echo "清空表时出错"
    exit 1
fi

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo "数据清空完成"
echo "耗时: $DURATION 秒"