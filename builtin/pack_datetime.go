package builtin

import (
	"fmt"
	"time"
)

// getTZ extracts the timezone from the env map. Returns UTC if not found.
func getTZ(env any) *time.Location {
	if m, ok := env.(map[string]any); ok {
		if tz, ok := m["_tz"]; ok {
			if loc, ok := tz.(*time.Location); ok {
				return loc
			}
		}
	}
	return time.UTC
}

// DATE creates a date from year, month, day components.
func DATE(args ...any) (any, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("DATE expects 3 arguments, got %d", len(args))
	}
	if args[0] == nil || args[1] == nil || args[2] == nil {
		return nil, nil
	}
	y, yNil := CoerceToFloat64(args[0])
	m, mNil := CoerceToFloat64(args[1])
	d, dNil := CoerceToFloat64(args[2])
	if yNil || mNil || dNil {
		return nil, nil
	}
	return time.Date(int(y), time.Month(int(m)), int(d), 0, 0, 0, 0, time.UTC), nil
}

// DATEVALUE converts a value to a date (time portion stripped).
func DATEVALUE(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("DATEVALUE expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()), nil
}

// DATETIMEVALUE converts a value to a datetime.
func DATETIMEVALUE(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("DATETIMEVALUE expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t, nil
}

// DAY returns the day of the month from a date.
func DAY(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("DAY expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t.Day(), nil
}

// MONTH returns the month from a date (1-12).
func MONTH(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("MONTH expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return int(t.Month()), nil
}

// YEAR returns the year from a date.
func YEAR(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("YEAR expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t.Year(), nil
}

// WEEKDAY returns the day of the week (SF convention: 1=Sunday, 7=Saturday).
func WEEKDAY(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("WEEKDAY expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return int(t.Weekday()) + 1, nil
}

// ADDMONTHS adds n months to a date, clamping month-end overflow.
func ADDMONTHS(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ADDMONTHS expects 2 arguments, got %d", len(args))
	}
	if args[0] == nil || args[1] == nil {
		return nil, nil
	}
	d, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	nf, nNil := CoerceToFloat64(args[1])
	if nNil {
		return nil, nil
	}
	n := int(nf)

	y, m, day := d.Date()
	totalMonths := int(m) - 1 + n
	targetYear := y + totalMonths/12
	targetMonth := time.Month(totalMonths%12 + 1)
	if totalMonths%12 < 0 {
		targetMonth += 12
		targetYear--
	}
	lastDay := time.Date(targetYear, targetMonth+1, 0, 0, 0, 0, 0, d.Location()).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(targetYear, targetMonth, day, d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location()), nil
}

// HOUR returns the hour from a datetime (0-23).
func HOUR(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("HOUR expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t.Hour(), nil
}

// MINUTE returns the minute from a datetime (0-59).
func MINUTE(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("MINUTE expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t.Minute(), nil
}

// SECOND returns the second from a datetime (0-59).
func SECOND(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SECOND expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t.Second(), nil
}

// MILLISECOND returns the millisecond from a datetime (0-999).
func MILLISECOND(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("MILLISECOND expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t.Nanosecond() / 1_000_000, nil
}

// TIMEVALUE returns the duration since midnight for a datetime.
func TIMEVALUE(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("TIMEVALUE expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	dt, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	y, m, d := dt.Date()
	midnight := time.Date(y, m, d, 0, 0, 0, 0, dt.Location())
	return dt.Sub(midnight), nil
}

// DAYOFYEAR returns the day of the year (1-366).
func DAYOFYEAR(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("DAYOFYEAR expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t.YearDay(), nil
}

// ISOWEEK returns the ISO 8601 week number.
func ISOWEEK(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ISOWEEK expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	_, week := t.ISOWeek()
	return week, nil
}

// ISOYEAR returns the ISO 8601 year.
func ISOYEAR(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ISOYEAR expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	year, _ := t.ISOWeek()
	return year, nil
}

// FROMUNIXTIME converts a Unix timestamp to a UTC time.
func FROMUNIXTIME(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("FROMUNIXTIME expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	ts, isNil := CoerceToFloat64(args[0])
	if isNil {
		return nil, nil
	}
	return time.Unix(int64(ts), 0).UTC(), nil
}

// UNIXTIMESTAMP converts a datetime to a Unix timestamp.
func UNIXTIMESTAMP(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("UNIXTIMESTAMP expects 1 argument, got %d", len(args))
	}
	if args[0] == nil {
		return nil, nil
	}
	t, ok := CoerceToTime(args[0])
	if !ok {
		return nil, nil
	}
	return t.Unix(), nil
}

// FORMATDURATION formats the duration between two datetimes as "DD:HH:MM:SS".
func FORMATDURATION(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("FORMATDURATION expects 2 arguments, got %d", len(args))
	}
	if args[0] == nil || args[1] == nil {
		return nil, nil
	}
	startDT, ok1 := CoerceToTime(args[0])
	endDT, ok2 := CoerceToTime(args[1])
	if !ok1 || !ok2 {
		return nil, nil
	}
	diff := endDT.Sub(startDT)
	if diff < 0 {
		diff = -diff
	}
	totalSecs := int(diff.Seconds())
	days := totalSecs / 86400
	hours := (totalSecs % 86400) / 3600
	minutes := (totalSecs % 3600) / 60
	seconds := totalSecs % 60
	return fmt.Sprintf("%d:%02d:%02d:%02d", days, hours, minutes, seconds), nil
}

// TODAY returns today's date at midnight in the configured timezone.
func ctxTODAY(env any, args ...any) (any, error) {
	tz := getTZ(env)
	now := time.Now().In(tz)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz), nil
}

// NOW returns the current datetime in the configured timezone.
func ctxNOW(env any, args ...any) (any, error) {
	tz := getTZ(env)
	return time.Now().In(tz), nil
}

// TIMENOW returns the duration since midnight in the configured timezone.
func ctxTIMENOW(env any, args ...any) (any, error) {
	tz := getTZ(env)
	now := time.Now().In(tz)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz)
	return now.Sub(midnight), nil
}

// DateTimeFunctions returns all datetime pack functions.
func DateTimeFunctions() []*PackFunction {
	variadicType := PackFuncTypes(new(func(...any) any))
	noArgType := PackFuncTypes(new(func() any))
	return []*PackFunction{
		{Name: "DATE", Fn: DATE, Types: variadicType},
		{Name: "DATEVALUE", Fn: DATEVALUE, Types: variadicType},
		{Name: "DATETIMEVALUE", Fn: DATETIMEVALUE, Types: variadicType},
		{Name: "DAY", Fn: DAY, Types: variadicType},
		{Name: "MONTH", Fn: MONTH, Types: variadicType},
		{Name: "YEAR", Fn: YEAR, Types: variadicType},
		{Name: "WEEKDAY", Fn: WEEKDAY, Types: variadicType},
		{Name: "ADDMONTHS", Fn: ADDMONTHS, Types: variadicType},
		{Name: "HOUR", Fn: HOUR, Types: variadicType},
		{Name: "MINUTE", Fn: MINUTE, Types: variadicType},
		{Name: "SECOND", Fn: SECOND, Types: variadicType},
		{Name: "MILLISECOND", Fn: MILLISECOND, Types: variadicType},
		{Name: "TIMEVALUE", Fn: TIMEVALUE, Types: variadicType},
		{Name: "DAYOFYEAR", Fn: DAYOFYEAR, Types: variadicType},
		{Name: "ISOWEEK", Fn: ISOWEEK, Types: variadicType},
		{Name: "ISOYEAR", Fn: ISOYEAR, Types: variadicType},
		{Name: "FROMUNIXTIME", Fn: FROMUNIXTIME, Types: variadicType},
		{Name: "UNIXTIMESTAMP", Fn: UNIXTIMESTAMP, Types: variadicType},
		{Name: "FORMATDURATION", Fn: FORMATDURATION, Types: variadicType},
		{Name: "TODAY", CtxFn: ctxTODAY, Types: noArgType, IsCtx: true},
		{Name: "NOW", CtxFn: ctxNOW, Types: noArgType, IsCtx: true},
		{Name: "TIMENOW", CtxFn: ctxTIMENOW, Types: noArgType, IsCtx: true},
	}
}
