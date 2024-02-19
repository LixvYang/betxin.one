package timeof

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

var zoneStr = time.Now().Format("-0700")

// 支持以下格式
// 1. 时间戳(秒)
// 2. 时间戳(毫秒)
// 2. 20060102
// 3. 20060102/00:00
// 4. 2006-01-02
// 5. 2006-01-02/00:00
// 6. 20060102150405
// 7. 2006-01-02T15:04:05Z07:00 (RFC3339)
// 8. duration-ago, 如 5h-ago
func TimeOf(str string) (t time.Time, ok bool) {
	var err error
	if strings.HasPrefix(str, "20") {
		str2 := str + zoneStr
		try := func(layout string) bool {
			t, err = time.Parse(layout+"-0700", str2)
			return err == nil
		}
		if try("20060102") || try("20060102/15:04") || try("2006-01-02") ||
			try("2006-01-02/15:04") || try("20060102150405") {
			ok = true
			return
		} else {
			t, err = time.Parse(time.RFC3339, str)
			t = t.Local()
			ok = err == nil
			return
		}
	} else if strings.HasSuffix(str, "-ago") {
		str = strings.TrimSuffix(str, "-ago")
		dur, err := time.ParseDuration(str)
		if err != nil {
			ok = false
			return
		}
		return time.Now().Add(-dur), true

	} else {
		// then, try to parse the time as timestamp in seconds.
		// this have bugs for parsing time before 2001-09-09...
		if len(str) == 10 {
			ts, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				return time.Unix(ts, 0), true
			}
		} else if len(str) == 13 {
			ts, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				return time.UnixMilli(ts), true
			}
		}
		return
	}
}

type stime []time.Time

func (p stime) Len() int           { return len(p) }
func (p stime) Less(i, j int) bool { return p[i].Before(p[j]) }
func (p stime) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortTime(in []time.Time) {
	sort.Sort(stime(in))
}
