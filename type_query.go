package telescope

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Query @Bean
type Query struct {
	Connection string   `json:"connection,omitempty"`
	Bindings   []string `json:"bindings,omitempty"`
	Sql        string   `json:"sql,omitempty"`
	Time       string   `json:"time,omitempty"`
	Slow       bool     `json:"slow,omitempty"`
	File       string   `json:"file,omitempty"`
	Line       string   `json:"line,omitempty"`
	Hash       string   `json:"hash,omitempty"`
	Hostname   string   `json:"hostname"`
}

func (b *Query) Init() {
	b.Hostname, _ = os.Hostname()
}

func (b *Query) BindType() string {
	return "query"
}

// QuerySplit 切割标识, 这个标识以后的代码才是业务的
var QuerySplit = "/app/entity/"

func (b *Query) Handler(entry *logrus.Entry) (*entries, []tag) {
	if strings.Index(entry.Message, "telescope_") != -1 {
		return nil, nil
	}
	// 根据模型目录 /app/entity/ 定位业务调用
	stack := string(debug.Stack())
	arr := strings.Split(string(stack), "\n")
	status := 0
	for _, str := range arr {
		if status <= 1 {
			index := strings.Index(str, QuerySplit)
			if index != -1 {
				status++                   // 模型调用自身,2次后 再下一层就是业务代码
				b.Connection = str[index:] // 这里记录模型
			}
		} else if strings.Index(str, QuerySplit) == -1 {
			// 第一个非模型目录
			arr2 := strings.Split(str, "/")
			if len(arr2) >= 4 {
				for i := len(arr2) - 4; i < len(arr2); i++ {
					b.File = b.File + "/" + arr2[i]
				}
			} else {
				for i := 0; i < len(arr2); i++ {
					b.File = b.File + "/" + arr2[i]
				}
			}
			break
		}
	}

	b.Sql = entry.Message
	t, ok := entry.Data["t"]
	if ok {
		b.Time = fmt.Sprintf("%.2f", t)
	}

	if b.File == "" {
		// 如果找不到业务级别的文件, 那就可能不是模型触发的, 这里是可见查询mysql.log
		b.File, b.Line = GetStackCallFile(stack, "go-home-admin/home/bootstrap/services/logs/mysql")
	}

	id := uuid.NewV4().String()
	return &entries{
		Uuid:                 id,
		BatchId:              NewtelescopeHook().TelescopeUUID(),
		FamilyHash:           nil,
		ShouldDisplayOnIndex: 1,
		Type:                 b.BindType(),
		Content:              ToContent(b),
		CreatedAt:            time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
