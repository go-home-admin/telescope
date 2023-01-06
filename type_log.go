package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

// Log @Bean
type Log struct {
	Level   string                 `json:"level,omitempty"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context"`
}

func (b Log) BindType() string {
	return "log"
}

func (b Log) Handler(entry *logrus.Entry) (*entries, []tag) {
	if entry.Level <= logrus.ErrorLevel {
		defer func() {
			telescopeEntries, tags := NewException().ToSave(string(debug.Stack()), entry.Message)
			NewtelescopeHook().Save(telescopeEntries, tags)
		}()
	}

	b.Message = entry.Message
	b.Context = entry.Data
	b.Level = entry.Level.String()
	uuid := uuid.NewV4().String()
	return &entries{
			Uuid:                 uuid,
			BatchId:              NewtelescopeHook().TelescopeUUID(),
			FamilyHash:           nil,
			ShouldDisplayOnIndex: 1,
			Type:                 b.BindType(),
			Content:              ToContent(b),
			CreatedAt:            time.Now().Format("2006-01-02 15:04:05"),
		}, []tag{{
			EntryUuid: uuid,
			Tag:       b.Level,
		}}
}
