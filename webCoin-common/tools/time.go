package tools

import "time"

// 转为iso时间
func ISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func ToTimeString(mill int64) string {
	milli := time.UnixMilli(mill)
	return milli.Format("2006-01-02 15:04:05")
}
