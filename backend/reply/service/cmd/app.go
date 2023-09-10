package main

import (
	"github.com/yehong-z/Cygnus/reply/interface/pkg"
	_ "github.com/yehong-z/Cygnus/reply/service/pkg/service"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_APP/conf/dubbogo.yaml
func main() {
	pkg.Init()
	// 以 API 的形式来启动框架
	rootConfig := config.NewRootConfigBuilder().
		SetConfigCenter(config.NewConfigCenterConfigBuilder().
			SetProtocol("nacos").SetAddress("121.36.89.81:8848"). // 根据配置结构，设置配置中心
			SetDataID("zyh_simple_IM").                           // 设置配置ID
			SetGroup("service").
			Build()).
		Build()

	if err := rootConfig.Init(); err != nil {
		panic(err)
	}
	select {}
}
