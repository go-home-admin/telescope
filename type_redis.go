package telescope

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"strings"
	"time"
)

// Redis @Bean
type Redis struct {
	Connection string `json:"connection"`
	Command    string `json:"command"`
	Time       string `json:"time"`
	Hostname   string `json:"hostname"`
}

func (b Redis) BindType() string {
	return "redis"
}

// 切割标识, 这个标识以后的代码才是业务的
var RedisSplit = "github.com/sirupsen/logrus/"

func (b Redis) Handler(entry *logrus.Entry) (*entries, []tag) {
	file, line := GetStackCallFile(string(debug.Stack()), RedisSplit)

	b.Connection = file + ":" + line
	b.Command = entry.Message
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

// 根据关键字的下一行获取调用文件和行号
func GetStackCallFile(stack, split string) (string, string) {
	var file, line string
	stack = strings.ReplaceAll(stack, "\n\t", "@NT@")
	arr := strings.Split(stack, "\n")
	status := 0

	file = arr[0]
	for _, str := range arr {
		index := strings.Index(str, split)
		if status == 0 {
			if index != -1 {
				status = 1
			}
		} else if status == 1 {
			if index == -1 {
				arr2 := strings.Split(str, "@NT@")
				if len(arr2) != 2 {
					break
				}
				file, line = strStackFile(arr2[1])
				break
			}
		}
	}

	return file, line
}
