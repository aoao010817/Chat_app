package trace
import (
	"io"
	"bytes"
	"testing"
)
func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("Hello trace")
		if buf.String() != "Hello trace\n" {
			t.Errorf("'%s'という謝った文字列が検出されました", buf.String())
		}
	}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}