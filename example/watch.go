package example

import (
	"fmt"
	"strconv"
	"time"

	"github.com/phamvinhdat/horus"
)

func main() {
	var v int
	w := horus.New("localhost:8500", horus.Pair{
		Key: "REDIS_MAX_CLIENTS",
		Obj: &v,
		ConvertFn: func(b []byte, obj interface{}) error {
			num, err := strconv.Atoi(string(b))
			if err != nil {
				return err
			}
			_obj := obj.(*int)
			*_obj = num
			return nil
		},
	})
	err := w.Watch()
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println(v)
		time.Sleep(time.Second)
	}
}
