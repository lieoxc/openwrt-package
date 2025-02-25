module("luci.controller.iotService", package.seeall)
local http = require "luci.http"

function index()
	entry({"admin", "services", "iotService"}, alias("admin", "services", "iotService", "iotWeb"), translate("iotService"), 100).dependent = true

    entry({"admin", "services", "iotService", "iotWeb"},template("iotService/iotWeb")).leaf = true
    entry({"admin", "services", "iotService", "get_lan_ip"}, call("get_lan_ip"), nil).leaf = true
end


function get_lan_ip()
    local uci = luci.model.uci.cursor()
    local lan_ip = uci:get("network", "lan", "ipaddr") or "192.168.1.1"  -- 默认回退地址
    luci.http.prepare_content("application/json")
    luci.http.write_json({ lan_ip = lan_ip })
end