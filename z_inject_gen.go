// gen for home toolset
package telescope

import (
	providers "github.com/go-home-admin/home/bootstrap/providers"
	gorm "gorm.io/gorm"
)

var _ProvidersSingle *Providers
var _telescopeHookSingle *telescopeHook

func GetAllProvider() []interface{} {
	return []interface{}{
		NewProviders(),
		NewtelescopeHook(),
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
