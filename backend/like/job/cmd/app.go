package main

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_APP/conf/dubbogo.yaml
func main() {
	// 以 API 的形式来启动框架
	rootConfig := config.NewRootConfigBuilder().
		SetConfigCenter(config.NewConfigCenterConfigBuilder().
			SetProtocol("nacos").SetAddress("121.36.89.81:8848"). // 根据配置结构，设置配置中心
			SetDataID("zyh_cygnus_like_service").                 // 设置配置ID
			SetGroup("like").
			Build()).
		Build()

	if err := rootConfig.Init(); err != nil {
		panic(err)
	}
	select {}
}
