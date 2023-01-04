package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

// Job @Bean
type Job struct {
	Status     string      `json:"status"`
	Connection string      `json:"connection"`
	Queue      string      `json:"queue"`
	Name       string      `json:"name"`
	Tries      interface{} `json:"tries"`
	Timeout    interface{} `json:"timeout"`
	Data       struct {
		Job                 interface{}   `json:"job"`
		Connection          interface{}   `json:"connection"`
		Queue               interface{}   `json:"queue"`
		ChainConnection     interface{}   `json:"chainConnection"`
		ChainQueue          interface{}   `json:"chainQueue"`
		ChainCatchCallbacks interface{}   `json:"chainCatchCallbacks"`
		Delay               interface{}   `json:"delay"`
		AfterCommit         interface{}   `json:"afterCommit"`
		Middleware          []interface{} `json:"middleware"`
		Chained             []interface{} `json:"chained"`
	} `json:"data"`
	Hostname string `json:"hostname"`
}

func (b Job) BindType() string {
	return "dump"
}

func (b Job) Handler(entry *logrus.Entry) (*entries, []tag) {
	b.Queue = entry.Message
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
