package telescope

import "github.com/sirupsen/logrus"

// Batch @Bean
type Batch struct {
	Queue          string    `json:"queue"`
	Connection     string    `json:"connection"`
	AllowsFailures string    `json:"allowsFailures"`
	Payload        []Payload `json:"payload"`
}

type Payload struct {
	id            string
	name          string
	totalJobs     string
	pendingJobs   string
	processedJobs string
	progress      string
	failedJobs    string
	options       string
	createdAt     string
	cancelledAt   string
	finishedAt    string
}

func (b Batch) BindType() string {
	return "batch"
}

func (b Batch) Handler(entry *logrus.Entry) (entries, []tag) {
	return entries{}, nil
}
