package logic

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

var Info string = `
{
	"name": "hello_effe",
	"version": "0.1",
	"doc" : "Getting start with effe"
}
`

type Context struct {
	value int64
}

func Init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Start() Context {
	fmt.Println("Start new Context")
	return Context{1 + rand.Int63n(2)}
}

func Run(ctx Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Hello from Effe with logs:  %d", ctx.value)
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
	return nil
}

func Stop(ctx Context) { return }
