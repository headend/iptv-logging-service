package monitor_logging_services

import (
	"context"
	"fmt"
	monitor_loggingpb "github.com/headend/iptv-monitor-service/proto"
	"github.com/headend/iptv-monitor-service/model"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/status"
	"log"
	"time"
)


func (c *monitorServerService) Gets(ctx context.Context, in *monitorpb.MonitorGetAll) (*monitorpb.MonitorResponse, error) {
	var monitorList []*model.Monitor
	err := c.DB.Db.Find(&monitorList).Error
	if err != nil {
		return nil, err
	}
	var res []*monitorpb.Monitor
	for _, tmp := range monitorList	{
		monitor := ConvertModelToProtoType(tmp)
		res = append(res, &monitor)
	}
	result := monitorpb.MonitorResponse{
		Status: monitorpb.MonitorResponseStatus_SUCCESS,
		Monitors: res,
	}
	return &result, nil
}

func (c *monitorServerService) Get(ctx context.Context, in *monitorpb.MonitorFilter) (*monitorpb.MonitorResponse, error) {
	var monitors []*model.Monitor
	mcase := matchFilterCase(in)
	var err error
	switch mcase {
	case 1:
		err = c.DB.Db.Where("id = ?", in.Id).Find(&monitors).Error
	case 23:
		err = c.DB.Db.Where("agent_id = : AND profile_id = ?", in.AgentId, in.ProfileId).Find(&monitors).Error
	case 2:
		err = c.DB.Db.Where("agent_id = ?", in.AgentId).Find(&monitors).Error
	case 3:
		err = c.DB.Db.Where("profile_id = ?", in.ProfileId).Find(&monitors).Error
	default:
		log.Printf("Exeption on %#v", in)
		return nil,nil
	}

	if err != nil {
		return nil, err
	}
	var res [] *monitorpb.Monitor

	for _,tmp := range monitors {
		monitor := ConvertModelToProtoType(tmp)
		res = append(res, &monitor)
	}
	result := monitorpb.MonitorResponse{
		Status: monitorpb.MonitorResponseStatus_SUCCESS,
		Monitors: res,
	}
	return &result, nil
}

func (c *monitorServerService) Add(ctx context.Context, in *monitorpb.MonitorRequest) (*monitorpb.MonitorResponse, error) {
	log.Println("params in: %v", in)
	notFound := c.DB.Db.Where("agent_id = : AND profile_id = ?", in.AgentId, in.ProfileId).RecordNotFound()
	if !notFound {
		err := fmt.Errorf("Agent is available! %#v", in)
		log.Println(err)
		return &monitorpb.MonitorResponse{
			Status: monitorpb.MonitorResponseStatus_FAIL,
			Monitors: nil,
		}, status.Errorf(409, "Confilic %s", err.Error())
	}
	monitorModel := monitorpb.Monitor{
		Id:            in.Id,
		StatusSignal:  in.StatusSignal,
		StatusVideo:   in.StatusVideo,
		StatusAudio:   in.StatusAudio,
		SignalMonitor: in.SignalMonitor,
		VideoMonitor:  in.VideoMonitor,
		AudioMonitor:  in.AudioMonitor,
		IsEnable:      in.IsEnable,
		AgentId:       in.AgentId,
		ProfileId:     in.ProfileId,
		StatusId:      in.StatusId,
	}
	var monitorProto monitorpb.Monitor
	err := c.DB.Db.Create(&monitorModel).Updates(map[string]interface{}{"date_create": time.Now(), "date_update": time.Now()}).Scan(&monitorProto).Error
	if err == nil {
		var res []*monitorpb.Monitor
		res = append(res, &monitorProto)
		return &monitorpb.MonitorResponse{Status: monitorpb.MonitorResponseStatus_SUCCESS, Monitors:res}, nil
	} else {
		log.Print(err)
		return &monitorpb.MonitorResponse{Status: monitorpb.MonitorResponseStatus_FAIL}, status.Errorf(409, "Confilic %s", err.Error())
	}

}

func (c *monitorServerService) Update(ctx context.Context, in *monitorpb.MonitorRequest) (*monitorpb.MonitorResponse, error) {
	var oldMonitorModel model.Monitor
	var err error

	if in.Id != 0 {
		err = c.DB.Db.Where("id = ?", in.Id).First(&oldMonitorModel).Error
	} else {
		err = c.DB.Db.Where("agent_id = : AND profile_id = ?", in.AgentId, in.ProfileId).First(&oldMonitorModel).Error
	}
	if err != nil {
		// Return not found if DB response notfound err
		if gorm.IsRecordNotFoundError(err) {
			return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_FAIL}, status.Error(404, "Not found")
		}
		// Return err
		return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_FAIL}, status.Error(500, "Internal server error")
	}
	// Return not found if can not found agent
	if oldMonitorModel == (model.Monitor{}) {
		return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_FAIL}, status.Error(404, "Not found")
	}

	return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_FAIL}, status.Errorf(403, "Not support")
}

func (c *monitorServerService) Delete(ctx context.Context, in *monitorpb.MonitorDelete) (*monitorpb.MonitorResponse, error) {
	var monitorModel model.Monitor
	var err error

	err = c.DB.Db.Where("id = ?", in.Id).First(&monitorModel).Error
	if err != nil {
		// Return not found if DB response notfound err
		if gorm.IsRecordNotFoundError(err) {
			return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_FAIL}, status.Error(404, "Not found")
		}
		// Return err
		return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_FAIL}, status.Error(500, "Internal server error")
	}
	// Return not found if can not found agent
	if monitorModel == (model.Monitor{}) {
		return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_FAIL}, status.Error(404, "Not found")
	}
	err2 := c.DB.Db.Delete(&monitorModel).GetErrors()
	var res []* monitorpb.Monitor
	tmp := ConvertModelToProtoType(&monitorModel)
	res = append(res, &tmp)
	if len(err2) == 0 {
		return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_SUCCESS, Monitors:res}, nil
	} else {
		log.Print(err2)
		return &monitorpb.MonitorResponse{Status:monitorpb.MonitorResponseStatus_FAIL}, status.Error(500, "Internal server error")
	}
}

func ConvertModelToProtoType(tmp *model.Monitor) monitorpb.Monitor {
	monitor := monitorpb.Monitor{
		Id : tmp.Id,
		StatusSignal : tmp.Status_signal,
		StatusVideo : tmp.Status_video,
		StatusAudio : tmp.Status_audio,
		SignalMonitor : tmp.Signal_monitor,
		VideoMonitor : tmp.Video_monitor,
		AudioMonitor : tmp.Audio_monitor,
		IsEnable : tmp.Is_enable,
		DateUpdate : tmp.Date_update.String(),
		AgentId : tmp.Agent_id,
		ProfileId : tmp.Profile_id,
		StatusId : tmp.Status_id,
	}
	return monitor
}

func matchFilterCase(in *monitorpb.MonitorFilter) (uint8) {

	if in.Id != 0 {
		return 1
	}
	if in.ProfileId != 0 && in.ProfileId != 0 {
		return 23
	}
	if in.AgentId != 0 {
		return 2
	}
	if in.ProfileId != 0 {
		return 3
	}
	return 0
}
