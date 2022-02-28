package datetime

import "time"

// DateTime 时间包装类
type DateTime struct {
	time.Time
}

// Today 今天
func (t DateTime) Today() DateTime {
	year, month, day := t.Date()
	time := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return DateTime{time}
}

// Seconds 获取总秒数
func (t DateTime) Seconds() int64 {
	return t.Unix()
}

// AddDays 添加指定的天数
func (t DateTime) AddDays(days int) DateTime {
	time := t.AddDate(0, 0, days)
	return DateTime{time}
}
