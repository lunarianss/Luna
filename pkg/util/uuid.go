package util

import (
	"strings"

	"github.com/google/uuid"
)

func GetUUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func GetShortUUID(len int) string {
	return GetUUID()[0:len]
}
