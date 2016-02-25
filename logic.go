// filename: lambda.go
package logic

import "net/http"
import "fmt"
import "math/rand"
import "time"

type Context struct{
    value int64
}

func Init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func Start() Context {
    fmt.Println("Start new Context")
    return Context{1 + rand.Int63n(2)};
}

func Run(ctx Context, w http.ResponseWriter, r *http.Request) error {
    fmt.Fprintf(w, "Hello from Logic: %d", ctx.value)
    return nil
}

func Stop(ctx Context) {return }

