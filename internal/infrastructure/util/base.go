package util

import (
	"encoding/json"

	"github.com/lunarianss/Luna/infrastructure/log"
)

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func LogCompleteInfo(v any) {
	c, _ := json.MarshalIndent(v, "", " ")
	log.Infof("%s", string(c))
}
