local dispatcher = require "luci.dispatcher"
local uci = require "luci.model.uci".cursor()
local http = require "luci.http"

m = Map("modem", translate("Dial Config"))
m.redirect = dispatcher.build_url("admin", "network", "modem","dial_overview")

s = m:section(NamedSection, arg[1], "dial-config", "")
s.addremove = false
s.dynamic = false
s:tab("general", translate("General Settings"))
s:tab("advanced", translate("Advanced Settings"))
s:tab("BandSetting", translate("Wirless Band Settings")) -- 频段设置

--------general--------

-- 是否启用
enable = s:taboption("general", Flag, "enable", translate("Enable"))
enable.default = "0"
enable.rmempty = false


-- 备注
remarks = s:taboption("general", Value, "remarks", translate("Remarks"))
remarks.rmempty = true

-- 移动网络
-- network = s:taboption("general", Value, "network", translate("Mobile Network"))
network = s:taboption("general", ListValue, "network", translate("Mobile Network"))
-- network.default = ""
network.rmempty = false

-- 获取移动网络，并显示设备名
function getMobileNetwork()
	local modem_number=uci:get('modem','global','modem_number')
	if modem_number == "0" then
		network:value("",translate("Mobile network not found"))
	end

	return
	uci:foreach("modem", "modem-device", function (modem_device)
		--获取模块名
		local modem_name=modem_device["name"]
		--获取网络
		local modem_network=modem_device["network"]

		--设置网络值
		network:value(modem_network,modem_network.." ("..translate(modem_name:upper())..")")
	end)
end

getMobileNetwork()

-- 拨号工具
dial_tool = s:taboption("advanced", ListValue, "dial_tool", translate("Dial Tool"))
dial_tool.description = translate("After switching the dialing tool, it may be necessary to restart the module or restart the router to recognize the module.")
dial_tool.default = "quectel-CM"
dial_tool.rmempty = true
-- dial_tool:value("", translate("Auto Choose"))
dial_tool:value("quectel-CM", translate("quectel-CM"))
-- dial_tool:value("mmcli", translate("mmcli"))

-- 网络类型
pdp_type= s:taboption("advanced", ListValue, "pdp_type", translate("PDP Type"))
pdp_type.default = "ipv4v6"
pdp_type.rmempty = false
pdp_type:value("ipv4", translate("IPv4"))
pdp_type:value("ipv6", translate("IPv6"))
pdp_type:value("ipv4v6", translate("IPv4/IPv6"))

-- 网络桥接
network_bridge = s:taboption("advanced", Flag, "network_bridge", translate("Network Bridge"))
network_bridge.description = translate("After checking, enable network interface bridge.")
network_bridge.default = "0"
network_bridge.rmempty = false

-- 接入点
apn = s:taboption("advanced", Value, "apn", translate("APN"))
apn.default = ""
apn.rmempty = true
apn:value("", translate("Auto Choose"))
apn:value("cmnet", translate("China Mobile"))
apn:value("3gnet", translate("China Unicom"))
apn:value("ctnet", translate("China Telecom"))
apn:value("cbnet", translate("China Broadcast"))
apn:value("5gscuiot", translate("Skytone"))

auth = s:taboption("advanced", ListValue, "auth", translate("Authentication Type"))
auth.default = "none"
auth.rmempty = false
auth:value("none", translate("NONE"))
auth:value("both", translate("PAP/CHAP (both)"))
auth:value("pap", "PAP")
auth:value("chap", "CHAP")

username = s:taboption("advanced", Value, "username", translate("PAP/CHAP Username"))
username.rmempty = true
username:depends("auth", "both")
username:depends("auth", "pap")
username:depends("auth", "chap")

password = s:taboption("advanced", Value, "password", translate("PAP/CHAP Password"))
password.rmempty = true
password.password = true
password:depends("auth", "both")
password:depends("auth", "pap")
password:depends("auth", "chap")

-- 配置ID
id = s:taboption("advanced", ListValue, "id", translate("Config ID"))
id.rmempty = false
id:value(arg[1])

-- 隐藏配置ID
m:append(Template("modem/hide_dial_config_id"))

-- 网络模式选择 4G，5G-SA , 5G-NSA
networkMode = s:taboption("BandSetting", ListValue, "networkMode",translate("Network Preferences"))
networkMode.default = "4G"
networkMode.rmempty = false
networkMode:value("4G",translate("4G"))
networkMode:value("5G-NSA",translate("5G-NSA"))
networkMode:value("5G-SA",translate("5G-SA"))

-- 频段锁定
band4 = s:taboption("BandSetting", MultiValue, "band4",translate("Band"))
band4.rmempty = true
band4.default = "auto"
band4:depends("networkMode", "4G")
band4:value("auto",translate("Auto Choose"))
band4:value("1",translate("B1"))
band4:value("2",translate("B2"))
band4:value("3",translate("B3"))
band4:value("5",translate("B5"))
band4:value("7",translate("B7"))
band4:value("8",translate("B8"))
band4:value("20",translate("B20"))
band4:value("38",translate("B28"))
band4:value("34",translate("B34"))
band4:value("38",translate("B38"))
band4:value("39",translate("B39"))
band4:value("40",translate("B40"))
band4:value("41",translate("B41"))

-- 5G 频段
band5 = s:taboption("BandSetting", MultiValue, "band5",translate("Band"))
band5.rmempty = true
band5.default = "auto"
band5:depends("networkMode", "5G-NSA")
band5:depends("networkMode", "5G-SA")
band5:value("auto",translate("Auto Choose"))
band5:value("1",translate("n1"))
band5:value("28",translate("n28"))
band5:value("41",translate("n41"))
band5:value("77",translate("n77"))
band5:value("78",translate("n78"))
band5:value("79",translate("n79"))


-- 频点
freq = s:taboption("BandSetting", Value, "freq",translate("Freq"))
freq.rmempty = true

-- PCI
pci = s:taboption("BandSetting", Value, "pci",translate("PCI"))
pci.rmempty = true

return m
