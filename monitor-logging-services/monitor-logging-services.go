package monitor_logging_services

import (
	"fmt"
	"github.com/headend/iptv-logging-service/model"
	monitor_loggingpb "github.com/headend/iptv-logging-service/proto"
	"github.com/headend/share-module/configuration"
	database "github.com/headend/share-module/databases"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type monitorLoggingServerService struct {
	config *configuration.Conf
	DB *database.Database
}


func StartServer()  {
	var config configuration.Conf

	config.LoadConf()
	// Make Database connect
	var db database.Database
	db.Connect(&config)
	if db.Err != nil {
		log.Fatal("Cannot connect DB: " + db.Err.Error())
	}
	// auto migrate data
	err2 := db.Db.AutoMigrate(model.MonitorLogs{}).Error
	if err2 != nil {
		panic(err2)
	}
	listenerAdd := fmt.Sprintf("%s:%d", config.RPC.Logging.Host, config.RPC.Logging.Port)
	ln, err := net.Listen("tcp", listenerAdd)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rpcServer := grpc.NewServer(grpc.MaxSendMsgSize(2 * 1024 *1024), grpc.ConnectionTimeout(10 * time.Second), grpc.MaxConcurrentStreams(2048))
	monitor_loggingpb.RegisterMonitorLogsServiceServer(rpcServer, &monitorLoggingServerService{config: &config, DB: &db})
	if rpcServer == nil {
		log.Fatalf("failed to register server: %v", err)
	}
	if err := rpcServer.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}