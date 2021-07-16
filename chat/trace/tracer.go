package trace
import (
	"io"
	"fmt"
)

type Tracer interface { //ログを保持
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct{}
func (t *nilTracer) Trace(a ...interface{}) {}
func Off() Tracer {
	return &nilTracer{}
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}