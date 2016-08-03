package configStore

import (
	"encoding/json"

	"github.com/tucnak/store"
)

func init() {
	store.Init("c14-cli")
	store.Register("json", json.Marshal, json.Unmarshal)
}

const (
	rcfile    = "c14rc.json"
	cachefile = "c14cache.json"
)

func GetCache(data interface{}) (err error) {
	err = store.Load(cachefile, &data)
	return
}

func GetRC(data interface{}) (err error) {
	err = store.Load(rcfile, data)
	return
}

func SaveCache(data interface{}) (err error) {
	err = store.SaveWith(cachefile, data, json.Marshal)
	return
}

func SaveRC(data interface{}) (err error) {
	err = store.SaveWith(rcfile, data, json.Marshal)
	return
}
