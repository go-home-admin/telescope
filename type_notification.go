package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

// Notification @Bean
type Notification struct {
	Notification string      `json:"notification"`
	Queued       bool        `json:"queued"`
	Notifiable   string      `json:"notifiable"`
	Channel      string      `json:"channel"`
	Response     interface{} `json:"response"`
	Hostname     string      `json:"hostname"`
}

func (b Notification) BindType() string {
	return "notification"
}

func (b Notification) Handler(entry *logrus.Entry) (*entries, []tag) {
	return &entries{
		Uuid:                 uuid.NewV4().String(),
		BatchId:              NewtelescopeHook().TelescopeUUID(),
		FamilyHash:           nil,
		ShouldDisplayOnIndex: 1,
		Type:                 b.BindType(),
		Content:              ToContent(b),
		CreatedAt:            time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
