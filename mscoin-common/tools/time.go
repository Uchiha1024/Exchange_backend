package tools

import "time"

func ISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func ToMill(str string) int64 {
	parse, _ := time.Parse("2006-01-02 15:04:05", str)
	return parse.UnixMilli()
}
func ToTimeString(mill int64) string {
	milli := time.UnixMilli(mill)
	return milli.Format("2006-01-02 15:04:05")
}

func ZeroTime() int64 {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.UnixMilli()
}
