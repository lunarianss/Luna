package util

import "encoding/json"

func DeepCopyUsingJSON[T interface{}](src, dst T) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
