package telescope

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Providers @Bean
type Providers struct {
	Mysql  *gorm.DB `inject:"database, @config(telescope.connect, default)"`
	isOpen bool
}

func (t *Providers) Init() {
	t.isOpen = true
	t.SetLog()
}

func (t *Providers) IsEnable() bool {
	return t.isOpen
}

// SetLog 打开望远镜时候, 设置log
func (t *Providers) SetLog() {
	hook := NewtelescopeHook()
	hook.mysql = t.Mysql
	hostname, _ := os.Hostname()
	hook.hostname = "home-server@" + hostname
	logrus.AddHook(hook)
}

// @Bean
type telescopeHook struct {
	mysql     *gorm.DB
	CidToUUID sync.Map
	hostname  string
}

func (t *telescopeHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (t *telescopeHook) Fire(entry *logrus.Entry) error {
	m, ok := entry.Data["type"]
	if !ok {
		m = EntryTypeLOG
	}
	mtype := m.(string)
	var content map[string]interface{}
	switch mtype {
	case EntryTypeQUERY:
		content, ok = t.EntryTypeQUERY(entry)
		if !ok {
			return nil
		}
	case EntryTypeREQUEST:
		content, ok = t.EntryTypeREQUEST(entry)
		if !ok {
			return nil
		}
	case EntryTypeREDIS:
		content = t.EntryTypeREDIS(entry)
	case EntryTypeJOB:
		content = t.EntryTypeJob(entry)
	default:
		content = t.EntryTypeLOG(entry)
	}
	contentStr, err := json.Marshal(content)
	if err != nil {
		contentStr = []byte("无法格式化log to json")
	}
	id := uuid.NewV4().String()
	data := map[string]interface{}{
		"uuid":                    id,
		"batch_id":                t.TelescopeUUID(),
		"family_hash":             nil,
		"should_display_on_index": 1,
		"type":                    mtype,
		"content":                 string(contentStr),
		"created_at":              time.Now().Format("2006-01-02 15:04:05"),
	}

	res := t.mysql.Table("telescope_entries").Create(data)
	if res.Error == nil {
		t.CreateTag(id, mtype, content)
	}
	return nil
}

func (t *telescopeHook) EntryTypeLOG(entry *logrus.Entry) map[string]interface{} {
	if entry.Level <= logrus.Level(logrus.ErrorLevel) {
		entry.Data["debug"] = string(debug.Stack())
	}
	return map[string]interface{}{
		"level":    entry.Level,
		"message":  entry.Message,
		"context":  entry.Data,
		"hostname": t.hostname,
	}
}

func (t *telescopeHook) EntryTypeQUERY(entry *logrus.Entry) (map[string]interface{}, bool) {
	if strings.Index(entry.Message, "telescope_") != -1 {
		return nil, false
	}

	var file, line string
	if entry.HasCaller() {
		file = entry.Caller.File
		line = strconv.Itoa(entry.Caller.Line)
	}
	return map[string]interface{}{
		"connection": "Mysql",
		"bindings":   "",
		"sql":        entry.Message,
		"time":       "0",
		"slow":       false,
		"file":       file,
		"line":       line,
		"hash":       "",
		"hostname":   t.hostname,
	}, true
}

func (t *telescopeHook) EntryTypeREQUEST(entry *logrus.Entry) (map[string]interface{}, bool) {
	var ctx interface{}
	var res interface{}
	ctx = entry.Context
	ginCtx := ctx.(*gin.Context)
	res = ginCtx.Writer
	telescopeResp := res.(*TelescopeResponseWriter)

	var response interface{}
	responseJSON := map[string]interface{}{}
	err := json.Unmarshal(telescopeResp.Body.Bytes(), &responseJSON)
	if err != nil || len(responseJSON) == 0 {
		response = telescopeResp.Body.String()
	} else {
		response = responseJSON
	}

	// 原始请求数据
	payload := make(map[string]interface{})
	if ginCtx.Request.PostForm == nil {
		raw, ok := ginCtx.Get("raw")
		if ok {
			data := raw.([]byte)
			if err == nil {
				_ = json.Unmarshal(data, &payload)
			}
		}
	} else {
		for k, v := range ginCtx.Request.PostForm {
			payload[k] = v[0]
		}
	}
	start, _ := ginCtx.Get("start")
	duration := time.Now().Sub(start.(time.Time))
	return map[string]interface{}{
		"ip_address": ginCtx.ClientIP(),
		"uri":        entry.Message,
		"method":     ginCtx.Request.Method,
		//"controller_action": "",
		//"middleware":        []string{},
		"headers": ginCtx.Request.Header,
		"payload": payload,
		//"session":           nil,
		"response_status": ginCtx.Writer.Status(),
		"response":        response,
		"duration":        duration.Milliseconds(),
		"memory":          ginCtx.Writer.Size(),
		"hostname":        t.hostname,
	}, true
}

func (t *telescopeHook) EntryTypeREDIS(entry *logrus.Entry) map[string]interface{} {
	return map[string]interface{}{
		"connection": "cache",
		"command":    entry.Message,
		"time":       "0",
		"hostname":   t.hostname,
	}
}

func (t *telescopeHook) EntryTypeJob(entry *logrus.Entry) map[string]interface{} {
	ginCtx := entry.Context.(*gin.Context)
	data, _ := ginCtx.Get("telescope_data")
	res := data.(map[string]interface{})
	res["hostname"] = t.hostname
	return res
}

func (t *telescopeHook) CreateTag(uuid, mType string, content map[string]interface{}) {
	var tag string
	switch mType {
	case "log":
		if _, ok := content["level"]; ok {
			tag = content["level"].(logrus.Level).String()
		}
	case "query":
		if _, ok := content["show"]; ok {
			if content["show"].(bool) {
				tag = "show"
			}
		}
	case "request":
		if _, ok := content["uri"]; ok {
			u, err := url.Parse(content["uri"].(string))
			if err == nil && u != nil {
				tag = u.Path
			}
		}
	case "job":
		if _, ok := content["status"]; ok {
			if content["status"].(string) == "failed" {
				tag = "failed"
			}
		}
	default:
		tag = ""
	}
	if tag != "" {
		t.mysql.Table("telescope_entries_tags").Create(map[string]interface{}{
			"entry_uuid": uuid,
			"tag":        tag,
		})
	}
}

func (t *telescopeHook) TelescopeUUID() string {
	cid := getGoId()
	v, ok := t.CidToUUID.Load(cid)

	if ok {
		return v.(string)
	}
	// 未开启情况就是使用，是无法关联的
	// 需要先使用 TelescopeStart()
	return time.Now().Format("2006-01-02 15:04:05")
}

func TelescopeStart() {
	cid := getGoId()
	t := NewtelescopeHook()
	t.CidToUUID.Store(cid, uuid.NewV4().String())
}

func TelescopeClose() {
	cid := getGoId()
	t := NewtelescopeHook()
	t.CidToUUID.Delete(cid)
}

// 获取跟踪ID, 严禁非开发模式使用
// github.com/bigwhite/experiments/blob/master/trace-function-call-chain/trace3/trace.go
func getGoId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

type TelescopeResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w TelescopeResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w TelescopeResponseWriter) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

var EntryTypeBATCH = "batch"
var EntryTypeCACHE = "cache"
var EntryTypeCOMMAND = "command"
var EntryTypeDUMP = "dump"
var EntryTypeEVENT = "event"
var EntryTypeEXCEPTION = "exception"
var EntryTypeJOB = "job"
var EntryTypeLOG = "log"
var EntryTypeMAIL = "mail"
var EntryTypeMODEL = "model"
var EntryTypeNOTIFICATION = "notification"
var EntryTypeQUERY = "query"
var EntryTypeREDIS = "redis"
var EntryTypeREQUEST = "request"
var EntryTypeSCHEDULED_TASK = "schedule"
var EntryTypeGATE = "gate"
var EntryTypeVIEW = "view"
var EntryTypeCLIENT_REQUEST = "home_request"
