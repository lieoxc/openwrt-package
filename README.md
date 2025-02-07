# lieoxc openwrt package

### luci-app-hypermodem
    4G/5G模组拨号Luci插件
### quectel_cm_5G
    移远模组拨号插件
### redis-patch
    官方openwrt的redis插件补丁，目的是修改官方插件的默认配置文件
    来源: https://github.com/openwrt/packages/tree/master/libs/redis
    改动： 内置一个redis.conf,编译安装的时候，自动安装到/etc/redis.conf
### mosquitto-patch
    官方openwrt的mosquitto插件补丁，目的是修改官方插件的默认配置文件
    来源: https://github.com/openwrt/packages/tree/master/net/mosquitto
    改动： 通过修改UCI配置,自动生成/etc/mosquitto/mosquitto.conf
### postgresql-patch
    官方openwrt的postgresql插件补丁，目的是修改官方插件的默认配置文件
    来源: https://github.com/openwrt/packages/tree/master/libs/postgresql
    改动：
        1. 修改UCI的配置：option PGDATA	/var/postgresql/data => 改为本地目录 /mnt/postgresql/data
        2. 允许外部链接，然后提前创建好对应的数据库以及用户名密码
            config postgresql config
                option PGDATA	/var/postgresql/data
            config postgres-db
                option name thingspanel
                option user postgres
                option pass postgresThingsPanel
        3. 允许外部访问（仅开发环境下需要）：
                1. 修改/mnt/postgresql/data/postgresql.conf 里面：
                    listen_addresses = '*' 监听所有接口地址
                2. 修改/mnt/postgresql/data/pg_hba.conf， 增加 ipv4栏目下：
                    host    all             all             0.0.0.0/0            trust

    
