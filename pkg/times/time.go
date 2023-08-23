package times

import "time"

// FormatTime 格式时间
func FormatTime(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

// GetYesterdayTimestamp 获取前一天时间戳
func GetYesterdayTimestamp() int64 {
	now := time.Now()
	return now.AddDate(0, 0, -1).Unix()
}
