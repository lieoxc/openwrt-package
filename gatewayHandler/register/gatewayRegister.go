package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gatewayhandler/model"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func gatewayRegister(devID string) error {
	addr := viper.GetString("web.addr")
	if addr == "" {
		addr = "localhost:9999"
		logrus.Println("Using default broker:", addr)
	}
	// 构造HTTP GET请求URL
	url := fmt.Sprintf("http://%s/api/v1/device/gateway-register", addr)
	logrus.Println("Request URL:", url)

	// 构造请求数据
	reqData := model.GatewayRegisterReq{
		GatewayId: devID, // 将设备 ID 放入请求数据中
		TenantId:  model.DefaultTenantId,
		Model:     model.DefaultGatewayCfgName,
	}

	// 将请求数据编码为 JSON
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		logrus.Println("Failed to marshal request data:", err)
		return err
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Println("Failed to send POST request:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Println("Failed to read response body:", err)
		return err
	}
	logrus.Debug(string(body))
	var respData model.HTTPRes
	err = json.Unmarshal(body, &respData)
	if err != nil {
		logrus.Println("Failed to json.Unmarshal:", err)
		return err
	}
	logrus.Debugf("MqttClientId:%v", respData.Data)
	return nil
}
