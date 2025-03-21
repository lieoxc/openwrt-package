require("luci.sys")

m = Map("autoreboot", translate("Scheduled Reboot"))
m.description = translate("Scheduled reboot Setting")

s = m:section(TypedSection, "login")
s.addremove = false
s.anonymous = true

enable = s:option(Flag,"enable",translate("Enable"))
enable.rmempty = false
enable.default = 0

week = s:option(ListValue, "week", translate("Week Day"))
week:value(7, translate("Everyday"))
week:value(1, translate("Monday"))
week:value(2, translate("Tuesday"))
week:value(3, translate("Wednesday"))
week:value(4, translate("Thursday"))
week:value(5, translate("Friday"))
week:value(6, translate("Saturday"))
week:value(0, translate("Sunday"))
week.default = 0

hour = s:option(Value, "hour", translate("Hour"))
hour.datatype = "range(0,23)"
hour.rmempty = false

pass = s:option(Value, "minute", translate("Minute"))
pass.datatype = "range(0,59)"
pass.rmempty = false

-- local e = luci.http.formvalue("cbi.apply")
-- if e then
--     -- 使用 luci.sys.call 同步执行命令
--     local result = luci.sys.call("/etc/init.d/autoreboot restart")
--     if result ~= 0 then
--         m.message = translate("Failed to apply configuration")
--     end
-- end

return m
