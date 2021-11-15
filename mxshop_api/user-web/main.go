package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/initialize"
	"mxshop_api/user-web/utils"
	myvalidator "mxshop_api/user-web/validator"

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

	viper.AutomaticEnv()
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

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
	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		2. 日志是分级别的，debug， info ， warn， error， fetal
		3. S函数和L函数很有用， 提供了一个全局的安全访问logger的途径
	*/
	zap.S().Infof("启动服务器,端口:%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", 8021)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
