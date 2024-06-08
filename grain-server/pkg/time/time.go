package timex

import (
	"fmt"
	"time"
)

const (
	YMDHMS = "2006-01-02 15:14:05"
	YMDHM  = "2006-01-02 15:14"
	YMD    = "2006-01-02"
	MY     = "2006-01"

	YMDHMSC = "2006年01年02 15:14:05"
	YMDHMC  = "2006年01年02 15:14"
	YMDC    = "2006年01年02"
	MYC     = "2006年01"
)

func GetSecondsLeftInDay() int64 {
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return int64(endOfDay.Sub(now).Seconds())
}

func GetStr2Time(timeStr string) (time.Time, error) {
	dob, err := time.Parse("2006-01-02 15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return dob, nil
}

func GetTime2Str(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func GetDateWithFormat(format string) string {
	return time.Now().Format(format)
}

func GetDate() string {
	return time.Now().Format("2006-01-02 15:04")
}

func GetDateYMD() string {
	return time.Now().Format("2006-01-02")
}

func GetTime2Date(t time.Time, format string) string {
	return t.Format(format)
}

func GetTimestampToDate(v int64, format string) string {
	tmp := time.UnixMilli(v)
	return tmp.Format(format)
}

func GetNowTimes() string {
	return time.Now().Format("20060102150405")
}

func GetDateBeforeHours(hour int) string {
	return time.Now().Add(-time.Hour * time.Duration(hour)).Format("2006-01-02 15:04:05")
}

// 获取今日之前的时间
func GetDateBeforeDays(days int) string {
	return time.Now().Add(-time.Hour * time.Duration(days) * 24).Format("2006-01-02")
}

// 获取今日之前的时间
func GetDateTimeBeforeDay(t time.Time, days int) string {
	return t.Add(-time.Hour * time.Duration(days) * 24).Format("2006-01-02")
}

// 获取今日之前的时间
func GetDateBeforeMinute(Minute int64) time.Time {
	now := time.Now()
	m, _ := time.ParseDuration(fmt.Sprintf("-%vm", Minute))
	m1 := now.Add(60 * m)
	return m1
}

func GetDateBeforeHour(Hour string) time.Time {
	now := time.Now()
	h, _ := time.ParseDuration(fmt.Sprintf("-%vh", Hour))
	h1 := now.Add(h)
	return h1
}

func GetAfterHour(Hour string) time.Time {
	now := time.Now()
	h, _ := time.ParseDuration(fmt.Sprintf("+%vh", Hour))
	h1 := now.Add(h)
	return h1
}

func GetHourToDuration(Hour string) time.Duration {
	h, _ := time.ParseDuration(fmt.Sprintf("+%vh", Hour))
	return h
}

func GetDateTimeBeforeDays(days int) string {
	return time.Now().Add(-time.Hour * time.Duration(days) * 24).Format("2006-01-02 15:04:05")
}

// 获取几天后
func GetDateAfterDays(days int) string {
	return time.Now().Add(+time.Hour * time.Duration(days) * 24).Format("2006-01-02")
}

func UnixMilliToDate(t int64) time.Time {
	tmp := time.UnixMilli(t)
	parse, _ := time.Parse(YMDHM, GetTime2Date(tmp, YMDHM))
	return parse
}

func GetStringMilliToDate(t string) time.Time {
	parse, _ := time.Parse("2006-01-02", t)
	return parse
}

func GetDateStringMilliToMilli(t string) int64 {
	parse, _ := time.Parse("2006-01-02", t)
	return parse.UnixMilli()
}

func GetStringToDate(t string, format string) time.Time {
	parse, _ := time.Parse(format, t)
	return parse
}

func GetMilliToTimeDate(seconds int64) string {
	day := seconds / 3600 / 24
	hour := seconds / 3600 % 24
	minute := seconds % 3600 / 60
	second := seconds % 60
	if day != 0 {
		return fmt.Sprintf("%v天%v小时%v分%v秒", day, hour, minute, second)
	}
	if hour != 0 {
		return fmt.Sprintf("%v小时%v分%v秒", hour, minute, second)
	}
	if minute != 0 {
		return fmt.Sprintf("%v分%v秒", minute, second)
	}
	if second != 0 {
		return fmt.Sprintf("%v秒", second)
	}
	return "xx分钟"
}
