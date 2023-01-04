// gen for home toolset
package telescope

import (
	providers "github.com/go-home-admin/home/bootstrap/providers"
	gorm "gorm.io/gorm"
	time "time"
)

var _ProvidersSingle *Providers
var _telescopeHookSingle *telescopeHook
var _BatchSingle *Batch
var _CacheSingle *Cache
var _ClientRequestSingle *ClientRequest
var _CommandSingle *Command
var _DumpSingle *Dump
var _EventSingle *Event
var _ExceptionSingle *Exception
var _JobSingle *Job
var _LogSingle *Log
var _ModelSingle *Model
var _QuerySingle *Query
var _RedisSingle *Redis
var _RequestSingle *Request
var _ScheduleSingle *Schedule
var _TcpSingle *Tcp

func GetAllProvider() []interface{} {
	return []interface{}{
		NewProviders(),
		NewtelescopeHook(),
		NewBatch(),
		NewCache(),
		NewClientRequest(),
		NewCommand(),
		NewDump(),
		NewEvent(),
		NewException(),
		NewJob(),
		NewLog(),
		NewModel(),
		NewQuery(),
		NewRedis(),
		NewRequest(),
		NewSchedule(),
		NewTcp(),
	}
}

func NewProviders() *Providers {
	if _ProvidersSingle == nil {
		_ProvidersSingle = &Providers{}
		_ProvidersSingle.Mysql = providers.GetBean("mysql").(providers.Bean).GetBean(*providers.GetBean("config").(providers.Bean).GetBean("telescope.connect, default").(*string)).(*gorm.DB)
		providers.AfterProvider(_ProvidersSingle, "")
	}
	return _ProvidersSingle
}
func NewtelescopeHook() *telescopeHook {
	if _telescopeHookSingle == nil {
		_telescopeHookSingle = &telescopeHook{}
		providers.AfterProvider(_telescopeHookSingle, "")
	}
	return _telescopeHookSingle
}
func NewBatch() *Batch {
	if _BatchSingle == nil {
		_BatchSingle = &Batch{}
		providers.AfterProvider(_BatchSingle, "")
	}
	return _BatchSingle
}
func NewCache() *Cache {
	if _CacheSingle == nil {
		_CacheSingle = &Cache{}
		providers.AfterProvider(_CacheSingle, "")
	}
	return _CacheSingle
}
func NewClientRequest() *ClientRequest {
	if _ClientRequestSingle == nil {
		_ClientRequestSingle = &ClientRequest{}
		providers.AfterProvider(_ClientRequestSingle, "")
	}
	return _ClientRequestSingle
}
func NewCommand() *Command {
	if _CommandSingle == nil {
		_CommandSingle = &Command{}
		providers.AfterProvider(_CommandSingle, "")
	}
	return _CommandSingle
}
func NewDump() *Dump {
	if _DumpSingle == nil {
		_DumpSingle = &Dump{}
		providers.AfterProvider(_DumpSingle, "")
	}
	return _DumpSingle
}
func NewEvent() *Event {
	if _EventSingle == nil {
		_EventSingle = &Event{}
		providers.AfterProvider(_EventSingle, "")
	}
	return _EventSingle
}
func NewException() *Exception {
	if _ExceptionSingle == nil {
		_ExceptionSingle = &Exception{}
		providers.AfterProvider(_ExceptionSingle, "")
	}
	return _ExceptionSingle
}
func NewJob() *Job {
	if _JobSingle == nil {
		_JobSingle = &Job{}
		providers.AfterProvider(_JobSingle, "")
	}
	return _JobSingle
}
func NewLog() *Log {
	if _LogSingle == nil {
		_LogSingle = &Log{}
		providers.AfterProvider(_LogSingle, "")
	}
	return _LogSingle
}
func NewModel() *Model {
	if _ModelSingle == nil {
		_ModelSingle = &Model{}
		providers.AfterProvider(_ModelSingle, "")
	}
	return _ModelSingle
}
func NewQuery() *Query {
	if _QuerySingle == nil {
		_QuerySingle = &Query{}
		providers.AfterProvider(_QuerySingle, "")
	}
	return _QuerySingle
}
func NewRedis() *Redis {
	if _RedisSingle == nil {
		_RedisSingle = &Redis{}
		providers.AfterProvider(_RedisSingle, "")
	}
	return _RedisSingle
}
func NewRequest() *Request {
	if _RequestSingle == nil {
		_RequestSingle = &Request{}
		providers.AfterProvider(_RequestSingle, "")
	}
	return _RequestSingle
}
func NewSchedule() *Schedule {
	if _ScheduleSingle == nil {
		_ScheduleSingle = &Schedule{}
		providers.AfterProvider(_ScheduleSingle, "")
	}
	return _ScheduleSingle
}
func NewTcp() *Tcp {
	if _TcpSingle == nil {
		_TcpSingle = &Tcp{}
		providers.AfterProvider(_TcpSingle, "")
	}
	return _TcpSingle
}
