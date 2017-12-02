package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func InterfaceToString(i interface{}) string {
	switch v := i.(type) {
	case time.Time:
		return v.In(time.Now().Location()).String()
	case bool:
		return strconv.FormatBool(i.(bool))
	case int:
		return strconv.FormatInt(int64(i.(int)), 10)
	case map[string]string:
		data, _ := json.Marshal(v)
		return string(data)
	default:
		return fmt.Sprintf("%s", v)
	}
}
