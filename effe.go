package main

import (
	"flag"
	"fmt"
	"github.com/siscia/effe/logic"
	"log/syslog"
	"net/http"
	"sync"
)

type complexContext struct {
	ctx logic.Context
	err error
}

func generateHandler(pool *sync.Pool, logger *syslog.Writer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := pool.Get().(complexContext)
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Crit("Logic Panicked")
			}
		}()
		err := logic.Run(ctx.ctx, ctx.err, w, r)
		if err != nil {
			logger.Debug(err.Error())
		}
		if ctx.err == nil {
			pool.Put(ctx)
		}
	}
}

func main() {
	port := flag.Int("port", 8080, "Port where serve the effe.")
	info := flag.Bool("info", false, "Print the effe information, then exit.")
	flag.Parse()
	if *info {
		fmt.Println(logic.Info)
		return
	}
	url := fmt.Sprintf(":%d", *port)
	logic.Init()
	logger, _ := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, "Logs From Effe ")
	var ctxPool = &sync.Pool{New: func() interface{} {
		ctx, err := logic.Start()
		return complexContext{ctx, err}
	}}
	http.HandleFunc("/", generateHandler(ctxPool, logger))
	http.ListenAndServe(url, nil)
}
