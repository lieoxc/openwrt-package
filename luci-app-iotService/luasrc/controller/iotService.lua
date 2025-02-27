module("luci.controller.iotService", package.seeall)
local http = require "luci.http"

function index()
    entry({"admin", "iotService"}, firstchild(), translate("IOT服务"), 25).dependent=false
	entry({"admin", "iotService", "iotService"}, alias("admin", "iotService", "iotService", "iotWeb"), translate("Device Manage Services"), 100).dependent = true

    entry({"admin", "iotService", "iotService", "iotWeb"},template("iotService/iotWeb")).leaf = true
    entry({"admin", "iotService", "iot_connect"},cbi("iotService/iot_connect"), translate("IOTConnect")).leaf = true
    entry({"admin", "iotService", "iotService", "get_lan_ip"}, call("get_lan_ip"), nil).leaf = true
end


function get_lan_ip()
    local uci = luci.model.uci.cursor()
    local lan_ip = uci:get("network", "lan", "ipaddr") or "192.168.1.1"  -- 默认回退地址
    luci.http.prepare_content("application/json")
    luci.http.write_json({ lan_ip = lan_ip })
end