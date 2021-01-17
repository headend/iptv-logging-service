package monitor_logging_services

import (
	"context"
	monitor_loggingpb "github.com/headend/iptv-logging-service/proto"
	"github.com/headend/iptv-logging-service/model"
	"google.golang.org/grpc/status"
	"log"
	"time"
)


func (c *monitorLoggingServerService) Gets(ctx context.Context, in *monitor_loggingpb.MonitorLogsGetAll) (*monitor_loggingpb.MonitorLogsResponse, error) {
	var monitorLogList []*model.MonitorLogs
	err := c.DB.Db.Find(&monitorLogList).Error
	if err != nil {
		return nil, err
	}
	var res []*monitor_loggingpb.MonitorLogs
	for _, tmp := range monitorLogList{
		monitor := ConvertModelToProtoType(*tmp)
		res = append(res, &monitor)
	}
	result := monitor_loggingpb.MonitorLogsResponse{
		Status: monitor_loggingpb.MonitorLogsResponseStatus_SUCCESS,
		MonitorLogs: res,
	}
	return &result, nil
}

func (c *monitorLoggingServerService) Get(ctx context.Context, in *monitor_loggingpb.MonitorLogsFilter) (*monitor_loggingpb.MonitorLogsResponse, error) {
	var monitorLogs []*model.MonitorLogs
	mcase := matchFilterCase(in)
	var err error
	switch mcase {
	case 1:
		err = c.DB.Db.Where("id = ?", in.AgentId).Find(&monitorLogs).Error
	case 2:
		err = c.DB.Db.Where("agent_id = ?", in.AgentId).Find(&monitorLogs).Error
	case 3:
		err = c.DB.Db.Where("profile_id = ?", in.ProfileId).Find(&monitorLogs).Error
	case 4:
		err = c.DB.Db.Where("monitor_id = ?", in.MonitorId).Find(&monitorLogs).Error
	case 5:
		err = c.DB.Db.Where("channel_id = ?", in.ChannelId).Find(&monitorLogs).Error
	case 6:
		err = c.DB.Db.Where("multicast_ip = ?", in.MulticastIp).Find(&monitorLogs).Error
	default:
		log.Printf("Exeption on %#v", in)
		return nil,nil
	}

	if err != nil {
		return nil, err
	}
	var res [] *monitor_loggingpb.MonitorLogs

	for _,tmp := range monitorLogs {
		monitorLog := ConvertModelToProtoType(*tmp)
		res = append(res, &monitorLog)
	}
	result := monitor_loggingpb.MonitorLogsResponse{
		Status: monitor_loggingpb.MonitorLogsResponseStatus_SUCCESS,
		MonitorLogs: res,
	}
	return &result, nil
}

func (c *monitorLoggingServerService) Add(ctx context.Context, in *monitor_loggingpb.MonitorLogsRequest) (*monitor_loggingpb.MonitorLogsResponse, error) {
	//log.Println("params in: %v", in)
	monitorLogModel := model.MonitorLogs{
		AgentID:            in.AgentId,
		ProfileID:          in.ProfileId,
		MonitorID:          in.MonitorId,
		ChannelID:          in.ChannelId,
		ChannelName:        in.ChannelName,
		MulticastIP:        in.MulticastIp,
		BeforeStatus:       in.BeforeStatus,
		BeforeSignalStatus: in.BeforeSignalStatus,
		BeforeVideoStatus:  in.BeforeVideoStatus,
		BeforeAudioStatus:  in.BeforeAudioStatus,
		AfterStatus:        in.AfterStatus,
		AfterSignalStatus:  in.AfterSignalStatus,
		AfterVideoStatus:   in.AfterVideoStatus,
		AfterAudioStatus:   in.AfterAudioStatus,
		Description:        in.Description,
	}
	var monitorLogProto monitor_loggingpb.MonitorLogs
	err := c.DB.Db.Create(&monitorLogModel).Updates(map[string]interface{}{"date_create": time.Now()}).Scan(&monitorLogProto).Error
	if err == nil {
		var res []*monitor_loggingpb.MonitorLogs
		res = append(res, &monitorLogProto)
		return &monitor_loggingpb.MonitorLogsResponse{Status: monitor_loggingpb.MonitorLogsResponseStatus_SUCCESS, MonitorLogs:res}, nil
	} else {
		log.Print(err)
		return &monitor_loggingpb.MonitorLogsResponse{Status: monitor_loggingpb.MonitorLogsResponseStatus_FAIL}, status.Errorf(409, "Confilic %s", err.Error())
	}

}

func matchFilterCase(in *monitor_loggingpb.MonitorLogsFilter) (uint8) {
	if in.Id != 0 {
		return 1
	}
	if in.AgentId != 0 {
		return 2
	}
	if in.ProfileId != 0 {
		return 3
	}
	if in.MonitorId != 0 {
		return 4
	}
	if in.ChannelId != 0 {
		return 5
	}
	if in.MulticastIp != "" {
		return 6
	}
	return 0
}

func ConvertModelToProtoType(logModel model.MonitorLogs) (logPb monitor_loggingpb.MonitorLogs) {
	result := monitor_loggingpb.MonitorLogs{
		Id:                 logModel.Id,
		AgentId:            logModel.AgentID,
		ProfileId:          logModel.ProfileID,
		MonitorId:          logModel.ProfileID,
		ChannelId:          logModel.ChannelID,
		ChannelName:        logModel.ChannelName,
		MulticastIp:        logModel.MulticastIP,
		BeforeStatus:       logModel.BeforeStatus,
		BeforeSignalStatus: logModel.BeforeSignalStatus,
		BeforeVideoStatus:  logModel.BeforeVideoStatus,
		BeforeAudioStatus:  logModel.BeforeAudioStatus,
		AfterStatus:        logModel.AfterStatus,
		AfterSignalStatus:  logModel.AfterSignalStatus,
		AfterVideoStatus:   logModel.AfterVideoStatus,
		AfterAudioStatus:   logModel.AfterAudioStatus,
		Description:        logModel.Description,
		DateCreate:         logModel.DateCreate.String(),
	}
	return result
}