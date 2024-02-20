package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

// Schedule @Bean
type Schedule struct {
	Command     string      `json:"command"`
	Description interface{} `json:"description"`
	Expression  string      `json:"expression"`
	Timezone    string      `json:"timezone"`
	User        interface{} `json:"user"`
	Output      string      `json:"output"`
	Hostname    string      `json:"hostname"`
}

func (s *Schedule) BindType() string {
	return "schedule"
}

func (s *Schedule) Handler(entry *logrus.Entry) (*entries, []tag) {
	b := *s

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
