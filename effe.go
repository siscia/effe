package main

import(
    "effe/logic"
    "net/http"
    "sync"
    "log/syslog"
    "flag"
    //"strconv"
    "fmt"
)

func generateHandler(pool *sync.Pool, logger *syslog.Writer) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request){
	ctx := pool.Get().(logic.Context)
	defer func() {
	    if r := recover(); r != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Crit("Logic Panicked")
	    }
	}()
	err := logic.Run(ctx, w, r)
	if err != nil {
	    logger.Debug(err.Error())
	}
	pool.Put(ctx)
    }
}

func main() {
    port := flag.Int("port", 8085, "Port where serve the effe.")
    flag.Parse()
    url := fmt.Sprintf(":%d", *port)
    logic.Init()
    logger, _ := syslog.New(syslog.LOG_ERR | syslog.LOG_USER, "Logs From Effe ")
    var ctxPool = &sync.Pool{New: func () interface{} {
	return logic.Start()} }
    http.HandleFunc("/", generateHandler(ctxPool, logger))
    http.ListenAndServe(url, nil)
}
