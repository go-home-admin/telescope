# Telescope
gin调试扩展, Telescope, 记录每个 API， 所有运行过程，包括请求参数、运行时间、执行 sql、打印 log、响应内容。


# 使用

默认使用`database.yaml`的`default`数据库。 自定义数据库，新增

`telescope.yaml`, 当前里面只有一行内容
````yaml
connection: mysql
````

注册中间件
````go
package http

import (
    "github.com/gin-gonic/gin"
    "github.com/go-home-admin/home/app"
    "github.com/go-home-admin/home/bootstrap/constraint"
    "github.com/go-home-admin/home/bootstrap/servers"
    "github.com/go-home-admin/telescope"
)

// Kernel @Bean
type Kernel struct {
    *servers.Http `inject:""`
}

func (k *Kernel) Init() {
    // 全局中间件
    k.Middleware = []gin.HandlerFunc{
        gin.Logger(),
        gin.Recovery(),
    }

    if app.IsDebug() {
        k.Middleware = append(k.Middleware, telescope.Providers())
    }

    // 分组中间件, 在路由提供者中自行设置
    k.MiddlewareGroup = map[string][]gin.HandlerFunc{
        "admin": {
            Cors(),
        },
        "api": {},
    }
}

// GetServer 提供统一命名规范的独立服务
func GetServer() constraint.KernelServer {
    return NewKernel()
}

````
## 查看调试信息
当前还不直接提供页面服务，只能在laravel框架下查看，[laravel](https://learnku.com/docs/laravel/9.x/installation/12200) 需要启动 [telescope](https://github.com/laravel/telescope) 。

请部署一个空的laravel即可，数据库要互通。
访问http://127.laravel.com/telescope/requests

