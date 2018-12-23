package core

import (
	"bytes"
	"math"
	"time"

	"github.com/d2r2/go-rsync/locale"
)

// FormatDurationToDaysHoursMinsSecs print time span
// in the format "x1 day(s) x2 hour(s) x3 minute(s) x4 second(s)".
// Understand plural cases for right spellings. Might be limited to number of sections.
func FormatDurationToDaysHoursMinsSecs(dur time.Duration, short bool, sections *int) string {
	var buf bytes.Buffer
	var totalHrs float64 = dur.Hours()
	days := totalHrs / 24
	count := 0
	if days >= 1 {
		count++
		var a int
		if sections == nil || *sections > count {
			a = int(days)
		} else {
			a = int(Round(days))
		}
		if short {
			buf.WriteString(f("%d %s", a, locale.TP(MsgDaysShort, nil, a)))
		} else {
			buf.WriteString(f("%d %s", a, locale.TP(MsgDaysLong, nil, a)))
		}
	}
	hours := totalHrs - float64(int(days)*24)
	if (hours >= 1 || count > 0) && (sections == nil || *sections > count) {
		if count > 0 {
			buf.WriteString(" ")
		}
		count++
		var a int
		if sections == nil || *sections > count {
			a = int(hours)
		} else {
			a = int(Round(hours))
		}
		if short {
			buf.WriteString(f("%d %s", a, locale.TP(MsgHoursShort, nil, a)))
		} else {
			buf.WriteString(f("%d %s", a, locale.TP(MsgHoursLong, nil, a)))
		}
	}
	var totalSecsLeft float64 = (dur - time.Duration(days)*24*time.Hour -
		time.Duration(hours)*time.Hour).Seconds()
	minutes := totalSecsLeft / 60
	if (minutes > 1 || count > 0) && (sections == nil || *sections > count) {
		if count > 0 {
			buf.WriteString(" ")
		}
		count++
		var a int
		if sections == nil || *sections > count {
			a = int(minutes)
		} else {
			a = int(Round(minutes))
		}
		if short {
			buf.WriteString(f("%d %s", a, locale.TP(MsgMinutesShort, nil, a)))
		} else {
			buf.WriteString(f("%d %s", a, locale.TP(MsgMinutesLong, nil, a)))
		}
	}
	seconds := int(totalSecsLeft - float64(int(minutes)*60))
	if (seconds > 0 || count == 0) && (sections == nil || *sections > count) {
		if count > 0 {
			buf.WriteString(" ")
		}
		if short {
			buf.WriteString(f("%d %s", seconds, locale.TP(MsgSecondsShort, nil, seconds)))
		} else {
			buf.WriteString(f("%d %s", seconds, locale.TP(MsgSecondsLong, nil, seconds)))
		}
	}
	return buf.String()
}

func pluralFloatToInt(val float64) int {
	if val == 1 {
		return 1
	} else if val < 1 {
		return 0
	} else if val < 2 {
		return 2
	} else {
		return int(Round(math.Floor(val)))
	}
}

const (
	kB = 1000
	MB = 1000 * 1000
	GB = 1000 * 1000 * 1000
	TB = 1000 * 1000 * 1000 * 1000
	PB = 1000 * 1000 * 1000 * 1000 * 1000
	EB = 1000 * 1000 * 1000 * 1000 * 1000 * 1000
)

func FormatSize(byteCount uint64, short bool) string {
	if byteCount > EB {
		a := float64(byteCount) / EB
		if short {
			return f("%v %s", a,
				locale.TP(MsgExaBytesShort, nil, pluralFloatToInt(a)))
		} else {
			return f("%v %s", a,
				locale.TP(MsgExaBytesLong, nil, pluralFloatToInt(a)))
		}
	} else if byteCount > PB {
		a := float64(byteCount) / PB
		if short {
			return f("%v %s", a,
				locale.TP(MsgPetaBytesShort, nil, pluralFloatToInt(a)))
		} else {
			return f("%v %s", a,
				locale.TP(MsgPetaBytesLong, nil, pluralFloatToInt(a)))
		}
	} else if byteCount > TB {
		a := float64(byteCount) / TB
		if short {
			return f("%v %s", a,
				locale.TP(MsgTeraBytesShort, nil, pluralFloatToInt(a)))
		} else {
			return f("%v %s", a,
				locale.TP(MsgTeraBytesLong, nil, pluralFloatToInt(a)))
		}
	} else if byteCount > GB {
		a := float64(byteCount) / GB
		if short {
			return f("%.1f %s", a,
				locale.TP(MsgGigaBytesShort, nil, pluralFloatToInt(a)))
		} else {
			return f("%.1f %s", a,
				locale.TP(MsgGigaBytesLong, nil, pluralFloatToInt(a)))
		}
	} else if byteCount > MB {
		a := int(Round(float64(byteCount) / MB))
		if short {
			return f("%v %s", a,
				locale.TP(MsgMegaBytesShort, nil, a))
		} else {
			return f("%v %s", a,
				locale.TP(MsgMegaBytesLong, nil, a))
		}
	} else if byteCount > kB {
		a := int(Round(float64(byteCount) / kB))
		if short {
			return f("%v %s", a,
				locale.TP(MsgKiloBytesShort, nil, a))
		} else {
			return f("%v %s", a,
				locale.TP(MsgKiloBytesLong, nil, a))
		}
	} else {
		a := int(byteCount)
		if short {
			return f("%v %s", a,
				locale.TP(MsgBytesShort, nil, a))
		} else {
			return f("%v %s", a,
				locale.TP(MsgBytesLong, nil, a))
		}
	}
}

func GetReadableSize(size FolderSize) string {
	return FormatSize(size.GetByteCount(), true)
}

func MegabytesToBytes(size uint64) uint64 {
	return size * MB
}
