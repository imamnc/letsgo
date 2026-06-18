package json

import stdjson "encoding/json"

func Marshal(value interface{}) (string, error) {
	b, err := stdjson.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Unmarshal[T any](data string) (T, error) {
	var value T
	if err := stdjson.Unmarshal([]byte(data), &value); err != nil {
		var zero T
		return zero, err
	}
	return value, nil
}
