package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"go.uber.org/zap"
	"mxshop_api/goods-web/global"
	"mxshop_api/goods-web/initialize"
	"mxshop_api/goods-web/utils/register/consul"
)

func main() {
	// 1.初始化logger
	initialize.InitLogger()
	// 2.初始化配置文件
	initialize.InitConfig()
	// 3.参数router
	Router := initialize.Routers()
	// 4.参数翻译器
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	// 5.初始化srv的链接
	initialize.InitSrvConn()

	register := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId, _ := uuid.NewV4()
	serviceIdStr := fmt.Sprintf("%s", serviceId)
	err := register.Register("localhost",
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		global.ServerConfig.Tags,
		serviceIdStr)

	if err != nil {
		zap.S().Infof("服务注册失败:%s", err.Error())
	}
	// 设置随机端口
	//viper.AutomaticEnv()
	//debug := viper.GetBool("MXSHOP_DEBUG")
	//if !debug {
	//	port, err := utils.GetFreePort()
	//	if err == nil {
	//		global.ServerConfig.Port = port
	//	}
	//}

	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		2. 日志是分级别的，debug， info ， warn， error， fetal
		3. S函数和L函数很有用， 提供了一个全局的安全访问logger的途径
	*/
	zap.S().Infof("启动服务器,端口:%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
