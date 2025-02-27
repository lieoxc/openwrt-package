package model

type GatewayRegisterReq struct {
	GatewayId string `json:"gateway_id"`
	TenantId  string `json:"tenant_id"`
	Model     string `json:"model"`
}

type GatewayRegisterRes struct {
	MqttUsername string `json:"mqtt_username"`
	MqttPassword string `json:"mqtt_password"`
	MqttClientId string `json:"mqtt_client_id"`
}
type HTTPRes struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	DefaultTenantId       = "d616bcbb"
	DefaultGatewayCfgName = "气象站"
	DefaultGatewayCfgID   = "d616bcbb"
)
