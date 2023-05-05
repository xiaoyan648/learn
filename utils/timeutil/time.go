package timeutil

import "time"

// 获取明天0点时刻时间
func GetTomorrowZeroTime(t time.Time) time.Time {
	y, m, d := time.Now().Date()
	zeroTimeToday := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	return zeroTimeToday.AddDate(0, 0, 1)
}
