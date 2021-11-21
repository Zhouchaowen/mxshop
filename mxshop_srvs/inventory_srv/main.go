package main

import (
	"flag"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/handler"
	"mxshop_srvs/inventory_srv/proto"
	"mxshop_srvs/inventory_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop_srvs/inventory_srv/initialize"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip addr")
	Port := flag.Int("port", 50051, "port")

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()

	zap.S().Info("ip: ", *IP)
	zap.S().Info("port: ", *Port)
	//if *Port == 0 {
	//	*Port,_ = utils.GetFreePort()
	//}

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen" + err.Error())
	}

	// 注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	register := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId, _ := uuid.NewV4()
	serviceIdStr := fmt.Sprintf("%s", serviceId)
	err = register.Register("localhost",
		*Port,
		global.ServerConfig.Name,
		global.ServerConfig.Tags,
		serviceIdStr)
	if err != nil {
		zap.S().Infof("服务注册失败:%s", err.Error())
	}
	zap.S().Debugf("启动服务器，端口：%d", *Port)

	// 启动
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register.DeRegister(serviceIdStr); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
