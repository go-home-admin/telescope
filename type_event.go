package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

// Event @Bean
type Event struct {
	Name    string `json:"name"`
	Payload struct {
		ConnectionName string `json:"connectionName"`
		Queue          string `json:"queue"`
		Payload        struct {
			Class      string `json:"class"`
			Properties struct {
				Value   string `json:"value"`
				Decoded struct {
					Uuid          string      `json:"uuid"`
					DisplayName   string      `json:"displayName"`
					Job           string      `json:"job"`
					MaxTries      interface{} `json:"maxTries"`
					MaxExceptions interface{} `json:"maxExceptions"`
					FailOnTimeout bool        `json:"failOnTimeout"`
					Backoff       interface{} `json:"backoff"`
					Timeout       interface{} `json:"timeout"`
					RetryUntil    interface{} `json:"retryUntil"`
					Data          struct {
						CommandName string `json:"commandName"`
						Command     string `json:"command"`
					} `json:"data"`
					TelescopeUuid string        `json:"telescope_uuid"`
					Id            string        `json:"id"`
					Attempts      int           `json:"attempts"`
					Type          string        `json:"type"`
					Tags          []interface{} `json:"tags"`
					PushedAt      string        `json:"pushedAt"`
				} `json:"decoded"`
			} `json:"properties"`
		} `json:"payload"`
	} `json:"payload"`
	Listeners []struct {
		Name   string `json:"name"`
		Queued bool   `json:"queued"`
	} `json:"listeners"`
	Broadcast bool   `json:"broadcast"`
	Hostname  string `json:"hostname"`
}

func (b Event) BindType() string {
	return "event"
}

func (b Event) Handler(entry *logrus.Entry) (*entries, []tag) {
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
