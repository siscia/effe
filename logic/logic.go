package logic

import (
	"fmt"
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

func Start() (Context, error) {
	fmt.Println("Start new Context")
	return Context{1 + rand.Int63n(2)}, nil
}

func Run(ctx Context, err error, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Hello from Effe:  %d\n", ctx.value)
	return nil
}

func Stop(ctx Context) { return }
