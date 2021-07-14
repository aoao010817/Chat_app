package trace

type Tracer interface { //ログを保持
	Trace(...interface{})
}