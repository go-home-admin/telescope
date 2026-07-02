package telescope

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

// Dump @Bean
type Dump struct {
	Dump string `json:"dump,omitempty"`
}

func (b Dump) BindType() string {
	return "dump"
}

func (b Dump) Handler(entry *logrus.Entry) (*entries, []tag) {
	b.Dump = entry.Message
	return &entries{
		Uuid:                 uuid.NewString(),
		BatchId:              NewtelescopeHook().TelescopeUUID(),
		FamilyHash:           nil,
		ShouldDisplayOnIndex: 1,
		Type:                 b.BindType(),
		Content:              ToContent(b),
		CreatedAt:            time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
