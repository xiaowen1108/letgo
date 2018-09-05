go API框架

```
    lg := letgo.New()
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