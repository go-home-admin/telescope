package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

// Redis @Bean
type Redis struct {
	Connection string `json:"connection"`
	Command    string `json:"command"`
	Time       string `json:"time"`
	Hostname   string `json:"hostname"`
}

func (b Redis) BindType() string {
	return "redis"
}

func (b Redis) Handler(entry *logrus.Entry) (*entries, []tag) {
	b.Command = entry.Message
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
