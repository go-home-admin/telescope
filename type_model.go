package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

// Model @Bean
type Model struct {
	Action   string `json:"action"`
	Model    string `json:"model"`
	Count    int    `json:"count"`
	Hostname string `json:"hostname"`
}

func (b Model) BindType() string {
	return "model"
}

func (b Model) Handler(entry *logrus.Entry) (*entries, []tag) {
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
