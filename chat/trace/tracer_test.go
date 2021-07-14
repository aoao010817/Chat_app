package trace
import (
	"bytes"
	"testing"
)
func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("")
		if buf.String() !=
	}
}