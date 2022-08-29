// gen for home toolset
package telescope

import (
	providers "github.com/go-home-admin/home/bootstrap/providers"
	gorm "gorm.io/gorm"
)

var _TelescopeProvidersSingle *Providers
var _telescopeHookSingle *telescopeHook

func GetAllProvider() []interface{} {
	return []interface{}{
		NewTelescopeProviders(),
		NewtelescopeHook(),
	}
}

func NewTelescopeProviders() *Providers {
	if _TelescopeProvidersSingle == nil {
		_TelescopeProvidersSingle = &Providers{}
		_TelescopeProvidersSingle.Mysql = providers.GetBean("database").(providers.Bean).GetBean(*(providers.GetBean("config").(providers.Bean).GetBean("telescope.connect").(*string))).(*gorm.DB)
		providers.AfterProvider(_TelescopeProvidersSingle, "")
	}
	return _TelescopeProvidersSingle
}
func NewtelescopeHook() *telescopeHook {
	if _telescopeHookSingle == nil {
		_telescopeHookSingle = &telescopeHook{}
		providers.AfterProvider(_telescopeHookSingle, "")
	}
	return _telescopeHookSingle
}
