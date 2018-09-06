go API框架

```
        lg := letgo.New()
    	lg.Error(func(cxt *letgo.Cxt) {
    		cxt.Send("傻逼方法不匹配哦")
    	}, func(cxt *letgo.Cxt) {
    		cxt.Send("傻逼找不到页面哦")
    	})
    	lg.Panic(func(writer http.ResponseWriter, request *http.Request, i interface{}) {
    		writer.Write([]byte("发生致命错误"))
    	})
    	lg.Get("/info/:userid", func(cxt *letgo.Cxt) {
    		cxt.Send(fmt.Sprintf("userid = %v", cxt.Get("userid")))
    	})
    	lg.Get("/", func(cxt *letgo.Cxt) {
    		cxt.Send("首页")
    	})
    	lg.Get("/aaa", func(cxt *letgo.Cxt) {
    		cxt.Send("aaa")
    	})
    	lg.Start(":8080")
```