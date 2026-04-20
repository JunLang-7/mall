package tools

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UUIDHex 生成一个唯一的 UUID，并将其中的连字符替换为空字符串，返回一个纯十六进制的 UUID 字符串
func UUIDHex() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func GetAllLike(str string) string {
	return "%" + str + "%"
}