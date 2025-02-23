package tools

import "time"

// 转为iso时间
func ISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
