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

func (l *Log) BindType() string {
	return "log"
}

func (l *Log) Handler(entry *logrus.Entry) (*entries, []tag) {
	b := *l

	if entry.Level <= logrus.ErrorLevel {
		hasError = true
		defer func() {
			telescopeEntries, tags := NewException().ToSave(string(debug.Stack()), entry.Message)
			NewtelescopeHook().Save(telescopeEntries, tags)
		}()
	}

	b.Message = entry.Message
	b.Context = entry.Data
	b.Level = entry.Level.String()
	id := uuid.NewV4().String()
	return &entries{
			Uuid:                 id,
			BatchId:              NewtelescopeHook().TelescopeUUID(),
			FamilyHash:           nil,
			ShouldDisplayOnIndex: 1,
			Type:                 b.BindType(),
			Content:              ToContent(b),
			CreatedAt:            time.Now().Format("2006-01-02 15:04:05"),
		}, []tag{{
			EntryUuid: id,
			Tag:       b.Level,
		}}
}
