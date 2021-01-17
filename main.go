package main

import (
	monitor_logging_service "github.com/headend/iptv-logging-service/monitor-logging-services"
	"log"
)


func main()  {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Logging service")
	monitor_logging_service.StartServer()
}
