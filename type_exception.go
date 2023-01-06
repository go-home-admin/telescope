package telescope

import (
	"bufio"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

// Exception @Bean
type Exception struct {
	Class       interface{}    `json:"class"`
	File        string         `json:"file"`
	Line        int            `json:"line"`
	Message     interface{}    `json:"message"`
	Context     interface{}    `json:"context"`
	Trace       []trace        `json:"trace"`
	LinePreview map[int]string `json:"line_preview"`
	Hostname    string         `json:"hostname"`
	Occurrences int            `json:"occurrences"`
}

func (b Exception) BindType() string {
	return "exception"
}

type trace struct {
	File string `json:"file"`
	Line string `json:"line"`
}

// 切割标识, 这个标识以后的代码才是业务的
var ExceptionSplit = "github.com/sirupsen/logrus/"

func (b Exception) Handler(entry *logrus.Entry) (*entries, []tag) {
	strStack := entry.Data["stack"].(string)
	return b.ToSave(strStack, entry.Message)
}

func (b Exception) ToSave(strStack, msg string) (*entries, []tag) {
	b.Message = msg
	b.Trace = make([]trace, 0)
	b.LinePreview = make(map[int]string)
	strStack = strings.ReplaceAll(strStack, "\n\t", "@NT@")
	arr := strings.Split(strStack, "\n")
	b.Class = arr[0]
	// 第一行是错误描述
	// 第二行开始, 代码函数@NT@文件:行
	// runtime/debug.Stack()@NT@/usr/local/opt/go/libexec/src/runtime/debug/stack.go:24 +0x7a
	for i := 1; i < len(arr); i++ {
		arr2 := strings.Split(arr[i], "@NT@")
		if len(arr2) != 2 {
			break
		}
		file, line := strStackFile(arr2[1])
		b.Trace = append(b.Trace, trace{
			File: file,
			Line: line,
		})
	}
	// 最后调用log的文件[github.com/sirupsen/logrus/exported.go] 那么下一行就是业务调用方
	inLogrus := 0
	for _, t := range b.Trace {
		if inLogrus == 0 {
			if strings.Index(t.File, ExceptionSplit) != -1 {
				inLogrus = 1
			}
		} else if inLogrus == 1 && strings.Index(t.File, ExceptionSplit) == -1 {
			b.File = t.File
			b.Line, _ = strconv.Atoi(t.Line)

			b.LinePreview = readFile(b.File, b.Line)
			break
		}
	}

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

// 分析文件和行号
// /Users/lvluo/Desktop/gongsi/nongchang.git/vendor/github.com/go-home-admin/telescope/type_log.go:24 +0x21a
func strStackFile(str string) (string, string) {
	space := strings.LastIndex(str, " ")
	if space <= 0 {
		return str, "0"
	}
	if len(str) < space {
		return str, "0"
	}
	str = str[:space]

	index := strings.LastIndex(str, ":")
	file := str[:index]
	line := str[index+1:]

	return file, line
}

// 获取文件指定行号的前后各十行内容
func readFile(path string, lineNumber int) map[int]string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	got := map[int]string{}
	fileScanner := bufio.NewScanner(file)
	lineCount := 1
	for fileScanner.Scan() {
		if lineCount > (lineNumber + 10) {
			break
		}

		if lineCount >= (lineNumber-20) && lineCount <= (lineNumber+10) {
			got[lineCount] = fileScanner.Text()
		}
		lineCount++
	}
	defer file.Close()
	return got
}
