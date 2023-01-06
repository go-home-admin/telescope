# Telescope
gin调试扩展, Telescope, 记录每个 API， 所有运行过程，包括请求参数、运行时间、执行 sql、打印 log、响应内容。


# 使用
当前还不直接提供页面服务，你要提前安装 [laravel](https://learnku.com/docs/laravel/9.x/installation/12200) 并且启用 [telescope](https://github.com/laravel/telescope) 。

go组件默认使用`database.yaml`的`default`数据库。 如需自定义数据库，新增 `telescope.yaml`, 当前里面只有一行内容。数据库要和laravel互通
````yaml
connect: mysql
````

#### 注册提供者

````go
package providers

import (
    "github.com/go-home-admin/home/bootstrap/constraint"
    "github.com/go-home-admin/home/bootstrap/providers"
    "github.com/go-home-admin/home/bootstrap/services"
    "github.com/go-home-admin/telescope"
)

// App @Bean
// 系统引导结构体
// 所有的服务提供者都应该在这里注入(注册)
type App struct {
    *services.Container          `inject:""`
    *providers.FrameworkProvider `inject:""`
    *providers.MysqlProvider     `inject:""`
    *providers.RedisProvider     `inject:""`

    *Route    `inject:""`
    *Response `inject:""`

    // 这是你需要加的代码，注册望远镜
    t *telescope.Providers `inject:""`
}

func (a *App) Run(servers []constraint.KernelServer) {
    a.Container.Run(servers)
}

````

#### 注册中间件
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
        k.Middleware = append(k.Middleware, telescope.Telescope())
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

````shell
docker run --name nw -p 8000:8000  feizhaoer/laravel-telescope:v2
docker exec -it nw /bin/bash
#docker: apt update
#docker: apt install vim
#docker: vim .env
````
访问http://127.0.0.1:8000/telescope/requests

## 非go-home框架下使用本组件

`telescope` 使用 `logrus` 的`Hook` 方式收集信息到数据库， 所以任意代码只要是使用`logrus.Info("msg)` 等等方式写的`log`。 
都可以在使用这个页面查看信息。

````go
// 初始化
t := &telescope.Providers{Mysql: gorm.DB}
t.Init()

// 记录任意log
logrus.WithFields(logrus.Fields{"type": "tcp", "read": raw, "status": 1, "start": start}).Debug("/route")
// 
````

