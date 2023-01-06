package telescope

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-home-admin/home/app"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Providers @Bean
type Providers struct {
	Mysql  *gorm.DB `inject:"mysql, @config(telescope.connect, default)"`
	isOpen bool

	init bool
}

var Routes = make(map[string]Type)

// SetDB 任意框架下使用， 需要手动设置DB
func (t *Providers) SetDB(db *gorm.DB) {
	t.Mysql = db
}

func (t *Providers) Init() {
	if app.IsDebug() && !t.init {
		t.init = true

		t.Register()
	}
}

func (t *Providers) Register() {
	for _, i := range GetAllProvider() {
		if v, ok := i.(Type); ok {
			Routes[v.BindType()] = v
		}
	}

	t.SetLog()
}

// SetLog 打开望远镜时候, 设置log
func (t *Providers) SetLog() {
	hook := NewtelescopeHook()
	hook.mysql = t.Mysql
	hostname, _ := os.Hostname()
	hook.hostname = "home-server@" + hostname
	logrus.AddHook(hook)
}

// AddRoute 允许重载处理
func (t *Providers) AddRoute(v Type) {
	Routes[v.BindType()] = v
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
		m = "log"
	}
	mType := m.(string)
	route, ok := Routes[mType]
	if ok {
		telescopeEntries, tags := route.Handler(entry)
		t.Save(telescopeEntries, tags)
	}

	return nil
}

func (t *telescopeHook) Save(telescopeEntries *entries, tags []tag) {
	if telescopeEntries != nil {
		res := t.mysql.Table("telescope_entries").Create(telescopeEntries)
		if res.Error == nil {
			for _, tag := range tags {
				t.mysql.Table("telescope_entries_tags").Create(map[string]interface{}{
					"entry_uuid": tag.EntryUuid,
					"tag":        tag.Tag,
				})
			}
		}
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
var EntryTypeTCP = "tcp"
var EntryTypeSCHEDULED_TASK = "schedule"
var EntryTypeGATE = "gate"
var EntryTypeVIEW = "view"
var EntryTypeCLIENT_REQUEST = "home_request"
