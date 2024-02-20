package telescope

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

// Tcp @Bean
type Tcp struct {
	IpAddress        string                 `json:"ip_address,omitempty"`
	Uri              string                 `json:"uri,omitempty"`
	Method           string                 `json:"method,omitempty"`
	ControllerAction string                 `json:"controller_action,omitempty"`
	Middleware       []string               `json:"middleware,omitempty"`
	Headers          map[string][]string    `json:"headers,omitempty"`
	Payload          map[string]interface{} `json:"payload,omitempty"`
	Session          map[string]interface{} `json:"session,omitempty"`
	ResponseStatus   int                    `json:"response_status,omitempty"`
	Response         interface{}            `json:"response,omitempty"`
	Duration         int64                  `json:"duration,omitempty"`
	Memory           int                    `json:"memory,omitempty"`
	Hostname         string                 `json:"hostname,omitempty"`
}

func (c *Tcp) BindType() string {
	return "tcp"
}

// Handler
// logrus.WithFields(logrus.Fields{"type": "tcp","read": raw,"tags": []string{str}}).Debug(tpc_route)
func (c *Tcp) Handler(entry *logrus.Entry) (*entries, []tag) {
	b := *c

	ip, ok := entry.Data["ip"]
	if ok {
		b.IpAddress = ip.(string)
	}
	// 原始请求数据
	raw, ok := entry.Data["read"]
	if ok {
		b.Payload = make(map[string]interface{})
		if ok {
			switch raw.(type) {
			case string:
				data := raw.(string)
				_ = json.Unmarshal([]byte(data), &b.Payload)
			case []byte:
				data := raw.([]byte)
				_ = json.Unmarshal(data, &b.Payload)
			}
		}
	}

	uuID := uuid.NewV4().String()
	b.Method = "TCP"
	b.Uri = entry.Message

	status, ok := entry.Data["status"]
	if ok {
		b.ResponseStatus = status.(int)
	}

	response, ok := entry.Data["response"]
	if ok {
		b.Response = response
	}

	controllerAction, ok := entry.Data["controller_action"]
	if ok {
		b.ControllerAction = controllerAction.(string)
	}
	start, ok := entry.Data["start"]
	if ok {
		b.Duration = time.Now().Sub(start.(time.Time)).Milliseconds()
	}

	retTags := make([]tag, 0)
	tags, ok := entry.Data["tags"]
	if ok {
		for _, t := range tags.([]string) {
			retTags = append(retTags, tag{
				EntryUuid: uuID,
				Tag:       t,
			})
		}
	}

	return &entries{
		Uuid:                 uuID,
		BatchId:              NewtelescopeHook().TelescopeUUID(),
		FamilyHash:           nil,
		ShouldDisplayOnIndex: 1,
		Type:                 "request",
		Content:              ToContent(b),
		CreatedAt:            time.Now().Format("2006-01-02 15:04:05"),
	}, retTags
}
