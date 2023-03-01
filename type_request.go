package telescope

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Request @Bean
type Request struct {
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

func (b Request) Init() {
	b.Hostname, _ = os.Hostname()
}

func (b Request) BindType() string {
	return "request"
}

func (b Request) Handler(entry *logrus.Entry) (*entries, []tag) {
	b.Payload = make(map[string]interface{})

	uuId := uuid.NewV4().String()
	var ctx interface{}
	var res interface{}
	ctx = entry.Context
	ginCtx := ctx.(*gin.Context)
	res = ginCtx.Writer
	telescopeResp := res.(*TelescopeResponseWriter)

	responseJSON := map[string]interface{}{}
	err := json.Unmarshal(telescopeResp.Body.Bytes(), &responseJSON)
	if err != nil || len(responseJSON) == 0 {
		b.Response = telescopeResp.Body.String()
	} else {
		b.Response = responseJSON
	}

	// 原始请求数据
	if ginCtx.Request.PostForm == nil {
		raw, ok := ginCtx.Get("raw")
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
	} else {
		for k, v := range ginCtx.Request.PostForm {
			b.Payload[k] = v[0]
		}
	}
	start, ok := ginCtx.Get("start")
	if ok {
		b.Duration = time.Now().Sub(start.(time.Time)).Milliseconds()
	}
	b.IpAddress = ginCtx.ClientIP()
	b.Uri = entry.Message
	b.Headers = ginCtx.Request.Header
	b.Method = ginCtx.Request.Method
	b.ResponseStatus = ginCtx.Writer.Status()

	return &entries{
		Uuid:                 uuId,
		BatchId:              NewtelescopeHook().TelescopeUUID(),
		FamilyHash:           nil,
		ShouldDisplayOnIndex: 1,
		Type:                 b.BindType(),
		Content:              ToContent(b),
		CreatedAt:            time.Now().Format("2006-01-02 15:04:05"),
	}, []tag{{Tag: ginCtx.FullPath(), EntryUuid: uuId}}
}
