package letgo

import (
	"net/http"
)

//顶层对象
type letgo struct {
	//middleware []MiddlewareFunc 中间件
	Route  *Router
	Server *http.Server
}

func (l *letgo) Start (port string)  {
	// 传入一个端口号，没有返回值，根据端口号开启http监听
	l.Server.Addr = port
	// 此处要进行资源初始化，加载所有路由、配置文件等等

	// 实例化env文件和config文件夹下的所有数据，根据配置

	// 根据路由列表，开始定义路由，并且根据端口号，开启http服务器
	l.Server.Handler = l.Route
	l.Server.ListenAndServe()
	// TODO 监听平滑升级和重启
}

func New() *letgo {
	return &letgo{Route:&Router{
		RedirectTrailingSlash:true,
		RedirectFixedPath:true,
		HandleMethodNotAllowed:true,
		HandleOPTIONS:true,
	}, Server:&http.Server{}}
}

//具体执行的逻辑处理单元
type Handle func(cxt *Cxt)
//Cxt
type Cxt struct {
	write http.ResponseWriter
	request *http.Request
	params []Param
}
func (cxt *Cxt) Send(data string) {
	cxt.write.Write([]byte(data))
}
func (cxt *Cxt) Get(key string) string {
	for _, v := range cxt.params {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}
//GET 方法
func (l *letgo) Get(path string, handler Handle) {
	l.Route.Handle("GET", path, handler)
}
//POST 方法
func (l *letgo) Post(path string, handler Handle) {
	l.Route.Handle("POST", path, handler)
}
//PUT 方法
func (l *letgo) Put(path string, handler Handle) {
	l.Route.Handle("PUT", path, handler)
}
//UPDATE 方法
func (l *letgo) Update(path string, handler Handle) {
	l.Route.Handle("UPDATE", path, handler)
}
//DELETE 方法
func (l *letgo) Delete(path string, handler Handle) {
	l.Route.Handle("DELETE", path, handler)
}

//Error 设置 MethodNotAllowed NotFound
func (l *letgo) Error(handlers... Handle) {
	for k, handle := range handlers {
		switch k {
			case 0:
				l.Route.MethodNotAllowed = handle
			case 1:
				l.Route.NotFound = handle
		}
	}
}
//Panic 设置 PanicHandler
func (l *letgo) Panic(handle func(http.ResponseWriter, *http.Request, interface{})) {
	l.Route.PanicHandler = handle
}