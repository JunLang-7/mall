package tools

import (
	"strings"

	"github.com/google/uuid"
)

// UUIDHex 生成一个唯一的 UUID，并将其中的连字符替换为空字符串，返回一个纯十六进制的 UUID 字符串
func UUIDHex() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
