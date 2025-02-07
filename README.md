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
