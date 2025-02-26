package register

// UUID 来源确定
// 每次启动都执行一遍
func DevRegisterInit() {
	gatewayRegister("001122334455")
	// 子设备的注册需要从 路由器的device表获取

}
