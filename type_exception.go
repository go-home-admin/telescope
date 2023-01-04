package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Exception @Bean
type Exception struct {
	Class       interface{} `json:"class"`
	File        interface{} `json:"file"`
	Line        interface{} `json:"line"`
	Message     interface{} `json:"message"`
	Context     interface{} `json:"context"`
	Trace       interface{} `json:"trace"`
	LinePreview []string    `json:"line_preview"`
	Hostname    string      `json:"hostname"`
	Occurrences int         `json:"occurrences"`
}

func (b Exception) BindType() string {
	return "exception"
}

type trace struct {
	file string
	line string
}

func (b Exception) Handler(entry *logrus.Entry) (*entries, []tag) {
	entry = entry.Data["entry"].(*logrus.Entry)
	stack := entry.Data["stack"].([]byte)

	linePreview := make([]string, 0)
	traces := make([]trace, 0)
	var files, lines, classs string
	status := -1
	arr := strings.Split(string(stack), "\n\t")
	for _, s := range arr {
		arr2 := strings.Split(s, "\n")
		var file, line, funs string
		switch len(arr2) {
		case 2:
			file = arr2[0]
			funs = arr2[1]
		case 1:
			file = arr2[0]
		}
		arrFile := strings.Split(file, ":")
		file = arrFile[0]
		line = arrFile[1]

		switch status {
		case -1:
			status = 0
			classs = s
		case 0:
			if strings.Count(s, "sirupsen/logrus") == 2 {
				status = 1 // 进入的log的代码
			}
		case 1:
			if strings.Count(s, "sirupsen/logrus") == 1 {
				status = 2 // 业务代码调用log
			}
		case 2:
			files = s
			lines = line
			status = 3
			traces = append(traces, trace{
				file: file,
				line: line,
			})
			linePreview = append(linePreview, funs)
		case 3:
			traces = append(traces, trace{
				file: file,
				line: line,
			})
			linePreview = append(linePreview, funs)
		}

	}
	b.Class = classs
	b.File = files
	b.Line = lines
	b.Message = entry.Message
	b.Context = entry.Data
	b.LinePreview = linePreview
	b.Occurrences = 1

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
