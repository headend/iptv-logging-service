package model

import "time"

func (MonitorLogs) TableName() string {
	return "monitor_logs"
}

type MonitorLogs struct {
	Id  int64 `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	AgentID int64 `gorm:"index:btree;column:agent_id;default:null"`
	ProfileID int64 `gorm:"index:btree;column:profile_id;default:null"`
	MonitorID  int64 `gorm:"index:btree;column:monitor_id;default:null"`
	ChannelID  int64 `gorm:"index:btree;column:channel_id;default:null"`
	ChannelName string `gorm:"index;column:channel_name;default:null"`
	MulticastIP string `gorm:"index;column:multicast_ip;default:null"`
	BeforeStatus int64 `gorm:"column:before_status;default:null"`
	BeforeSignalStatus  bool  `gorm:"column:before_signal_status;default:false"`
	BeforeVideoStatus bool  `gorm:"column:before_video_status;default:false"`
	BeforeAudioStatus bool  `gorm:"column:before_audio_status;default:false"`
	AfterStatus int64 `gorm:"column:after_status;default:null"`
	AfterSignalStatus  bool  `gorm:"column:after_signal_status;default:false"`
	AfterVideoStatus bool  `gorm:"column:after_video_status;default:false"`
	AfterAudioStatus bool  `gorm:"column:after_audio_status;default:false"`
	Description string `gorm:"column:desc;default:null"`
	DateCreate time.Time `gorm:"column:date_create;default:null"`
}