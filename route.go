package letgo

import (
	"net/http"
)
//路由对象
type Router struct {
	entry map[string]*node
	RedirectTrailingSlash bool
	RedirectFixedPath bool
	HandleMethodNotAllowed bool
	HandleOPTIONS bool
	MethodNotAllowed Handle
	NotFound Handle
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}

//路由节点
//path   路由path
//param  路由参数
//handle 具体执行的逻辑处理单元
//type node struct {
//	path   string
//	param  map[int]string
//	handle Handle
//	method string
//}
//转换path
//func transformation(path string) (string, map[int]string, error) {
//	path = strings.TrimLeft(path, "/")
//	ps := strings.Split(path, "/")
//	if len(ps) > 1 {
//		param := make(map[int]string)
//		for k, v := range ps {
//			if strings.HasPrefix(v, ":") {
//				param[k] = strings.TrimLeft(v, ":")
//			}
//		}
//		if len(param) > 0 {
//			return "/" + ps[0], param, nil
//		}
//	}
//	return "/" + ps[0], nil, nil
//}
//路由添加
func (r *Router) Handle(method string, path string, handler Handle) {
	//解析path
	//path, param, err := transformation(path)
	//if err != nil {
	//	log.Fatal("Error:", err)
	//}
	//n := &node{path:path, param:param, handle:handler, method:method}
	//r.entry[path] = n
	if path[0] != '/' {
		//panic("path must begin with '/' in path '" + path + "'")
		path = "/" + path
	}
	if r.entry == nil {
		r.entry = make(map[string]*node)
	}
	root, ok := r.entry[method]
	if !ok {
		root = new(node)
		r.entry[method] = root
	}
	root.addRoute(path, handler)
}
func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}
func (r *Router) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range r.entry {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range r.entry {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _, _ := r.entry[method].getValue(path)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}
func (r *Router) ServeHTTP (write http.ResponseWriter, request *http.Request)  {
	////解析请求
	//path := request.URL.Path
	////屏蔽   /favicon.ico
	//if path != "/favicon.ico" {
	//	//解析path
	//	path = strings.TrimLeft(path, "/")
	//	ps := strings.Split(path, "/")
	//	param := make(map[int]string)
	//	for k, v := range ps {
	//		if k == 0 {
	//			path = "/" + v
	//		} else {
	//			param[k] = v
	//		}
	//	}
	//	if node, ok := r.entry[path]; !ok {
	//		//未匹配  404
	//		http.NotFound(write, request)
	//	} else {
	//		if request.Method != node.method && node.method != "ANY"{
	//			http.NotFound(write, request)
	//		} else {
	//			_param := make(map[string]string)
	//			if node.param != nil {
	//				for k, v := range node.param {
	//					if _v, ok := param[k]; ok{
	//						_param[v] = _v
	//					}
	//				}
	//			}
	//			node.handle(&Cxt{write,request, _param})
	//		}
	//	}
	//}
	if r.PanicHandler != nil {
		//做一些异常处理操作
		defer r.recv(write, request)
	}
	path := request.URL.Path
	if root := r.entry[request.Method]; root != nil {
		if handle, ps, tsr := root.getValue(path); handle != nil {
			handle(&Cxt{write,request, ps})
			return
		} else if request.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if request.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}
			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					request.URL.Path = path[:len(path)-1]
				} else {
					request.URL.Path = path + "/"
				}
				http.Redirect(write, request, request.URL.String(), code)
				return
			}
			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					request.URL.Path = string(fixedPath)
					http.Redirect(write, request, request.URL.String(), code)
					return
				}
			}
		}
	}

	if request.Method == "OPTIONS" && r.HandleOPTIONS {
		// Handle OPTIONS requests
		if allow := r.allowed(path, request.Method); len(allow) > 0 {
			write.Header().Set("Allow", allow)
			return
		}
	} else {
		// Handle 405
		if r.HandleMethodNotAllowed {
			if allow := r.allowed(path, request.Method); len(allow) > 0 {
				write.Header().Set("Allow", allow)
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed(&Cxt{write,request, nil})
				} else {
					http.Error(write,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}
	// Handle 404
	if r.NotFound != nil {
		r.NotFound(&Cxt{write,request, nil})
	} else {
		http.NotFound(write, request)
	}
}

// CleanPath is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//	1. Replace multiple slashes with a single slash.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned
func CleanPath(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}

// internal helper to lazily create a buffer if necessary
func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}
