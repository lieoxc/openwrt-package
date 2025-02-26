mp = Map("iot_connect")
mp.title = translate("IOTConnect")
mp.description = translate("Set The Config To Connect IOT Platform")

s = mp:section(TypedSection, "service", translate("MQTT Setting"))
s.anonymous = true

publicEnabled = s:option(Flag, "publicEnabled", translate("Enable Public"))
publicEnabled.description = translate("After Enable, router will report data to Public Server.")
publicEnabled.default = 0
publicEnabled.rmempty = false

publicAddr = s:option(Value, "publicAddr", translate("publicAddr"))
publicAddr.default = ""
publicAddr.rmempty = true

publicPort = s:option(Value, "publicPort", translate("publicPort"))
publicPort.default = ""
publicPort.rmempty = true



privateEnabled = s:option(Flag, "privateEnabled", translate("Enable Private"))
privateEnabled.description = translate("After Enable, router will report data to Private Server.")
privateEnabled.default = 0
privateEnabled.rmempty = false

privateAddr = s:option(Value, "privateAddr", translate("privateAddr"))
privateAddr.default = ""
privateAddr.rmempty = true


privatePort = s:option(Value, "privatePort", translate("privatePort"))
privatePort.default = ""
privatePort.rmempty = true

return mp