package horus

type Watcher interface {
	Watch() error
}

type convertFn func(b []byte, obj interface{}) error

type Pair struct {
	Key       string
	Obj       interface{}
	ConvertFn convertFn
}
