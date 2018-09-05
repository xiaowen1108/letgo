package letgo

import (
	"net/http"
	"log"
	"strings"
)
//路由对象
type Router struct {
	entry map[string]*node
}

//路由节点
//path   路由path
//param  路由参数
//handle 具体执行的逻辑处理单元
type node struct {
	path   string
	param  map[int]string
	handle Handle
	method string
}
//转换path
func transformation(path string) (string, map[int]string, error) {
	path = strings.TrimLeft(path, "/")
	ps := strings.Split(path, "/")
	if len(ps) > 1 {
		param := make(map[int]string)
		for k, v := range ps {
			if strings.HasPrefix(v, ":") {
				param[k] = strings.TrimLeft(v, ":")
			}
		}
		if len(param) > 0 {
			return "/" + ps[0], param, nil
		}
	}
	return "/" + ps[0], nil, nil
}
//路由添加
func (r *Router) Handle(method string, path string, handler Handle) {
	//解析path
	path, param, err := transformation(path)
	if err != nil {
		log.Fatal("Error:", err)
	}
	n := &node{path:path, param:param, handle:handler, method:method}
	r.entry[path] = n
}

func (r *Router) ServeHTTP (write http.ResponseWriter, request *http.Request)  {
	//解析请求
	path := request.URL.Path
	//屏蔽   /favicon.ico
	if path != "/favicon.ico" {
		//解析path
		path = strings.TrimLeft(path, "/")
		ps := strings.Split(path, "/")
		param := make(map[int]string)
		for k, v := range ps {
			if k == 0 {
				path = "/" + v
			} else {
				param[k] = v
			}
		}
		if node, ok := r.entry[path]; !ok {
			//未匹配  404
			http.NotFound(write, request)
		} else {
			if request.Method != node.method && node.method != "ANY"{
				http.NotFound(write, request)
			} else {
				_param := make(map[string]string)
				if node.param != nil {
					for k, v := range node.param {
						if _v, ok := param[k]; ok{
							_param[v] = _v
						}
					}
				}
				node.handle(&Cxt{write,request, _param})
			}
		}
	}
}
