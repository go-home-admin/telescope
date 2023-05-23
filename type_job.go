package telescope

import (
	"fmt"
	app2 "github.com/go-home-admin/home/app"
	"github.com/go-home-admin/home/bootstrap/constraint"
	"github.com/go-home-admin/home/bootstrap/servers"
	"github.com/go-home-admin/home/bootstrap/services/app"
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
	Data       interface{} `json:"data"`
	Hostname   string      `json:"hostname"`
}

func (b *Job) Boot() {
	if app2.IsDebug() && app.HasBean("queue") {
		server := app.GetBean("queue").(*servers.Queue)
		server.AddMiddleware(func(job constraint.Job, next func(constraint.Job)) {
			TelescopeStart()
			defer TelescopeClose()
			defer func() {
				// 记录调试信息
				logrus.WithFields(logrus.Fields{
					"type":       "job",
					"status":     "processed",
					"connection": app2.Config("queue.connection", "redis"),
					"queue":      app2.Config("queue.queue.stream_name", "home_default_stream"),
					"data":       job,
				}).Debug(fmt.Sprintf("%T", job))
			}()

			next(job)
		})
	}
}

func (b *Job) BindType() string {
	return "job"
}

func (b *Job) Handler(entry *logrus.Entry) (*entries, []tag) {
	b.Name = entry.Message
	b.Data = entry.Data["data"]
	b.Status = entry.Data["status"].(string)
	b.Queue = entry.Data["queue"].(string)
	b.Connection = entry.Data["connection"].(string)

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
