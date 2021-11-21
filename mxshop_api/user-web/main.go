package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/initialize"
	"mxshop_api/user-web/utils/register/consul"
	myvalidator "mxshop_api/user-web/validator"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
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

	//viper.AutomaticEnv()
	//debug := viper.GetBool("MXSHOP_DEBUG")
	//if !debug {
	//	port, err := utils.GetFreePort()
	//	if err == nil {
	//		global.ServerConfig.Port = port
	//	}
	//}

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// 6.注册中心
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

	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		2. 日志是分级别的，debug， info ， warn， error， fetal
		3. S函数和L函数很有用， 提供了一个全局的安全访问logger的途径
	*/
	zap.S().Infof("启动服务器,端口:%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit

	err = register.DeRegister(serviceIdStr)
	if err != nil {
		zap.S().Panic("注销失败:", err.Error())
	}else {
		zap.S().Panic("注销成功")
	}
}
