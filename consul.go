package horus

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type consulWatcher struct {
	addr       string
	convertFns map[string]convertFn
	pair       map[string]interface{}
}

func (c *consulWatcher) Watch() error {
	plan, err := c.parsePlan()
	if err != nil {
		log.Println("failed to parse plan", "error", err)
		return err
	}

	plan.Handler = func(_ uint64, raw interface{}) {
		switch d := raw.(type) {
		case *api.KVPair:
			obj, ok := c.pair[d.Key]
			if !ok {
				return
			}

			converter := defaultConvertFn
			if objConverter, ok := c.convertFns[d.Key]; ok {
				converter = objConverter
			}
			err = converter(d.Value, obj)
			if err != nil {
				log.Println("failed convert object of watcher", "error", err)
			}
		default:
			log.Println("type of event unsupported")
			return
		}
	}

	go func() {
		err = plan.Run(c.addr)
		if err != nil {
			log.Panic("failed to sync", "error", err)
		}
	}()

	return nil
}

func New(addr string, pairs ...Pair) Watcher {
	if len(pairs) == 0 {
		log.Panic("pair is empty")
	}

	consul := consulWatcher{
		addr:       addr,
		convertFns: make(map[string]convertFn),
		pair:       make(map[string]interface{}),
	}

	for _, pair := range pairs {
		consul.pair[pair.Key] = pair.Obj
		if pair.ConvertFn != nil {
			consul.convertFns[pair.Key] = pair.ConvertFn
		}
	}

	return &consul
}

func (c *consulWatcher) parsePlan() (*watch.Plan, error) {
	params := map[string]interface{}{
		"type": "key",
	}
	for key, _ := range c.pair {
		params["key"] = key
	}
	return watch.Parse(params)
}

func defaultConvertFn(b []byte, obj interface{}) error {
	return json.Unmarshal(b, obj)
}
