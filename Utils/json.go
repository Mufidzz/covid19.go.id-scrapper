package Utils

import "encoding/json"

func Stringify(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(str)
}
