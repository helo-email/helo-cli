package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Timestamp and date flags take either an absolute value or a duration back
// from now, so a window like "the last day" needs no date arithmetic:
// --from 24h, --since 7d, --start 1d12h. Units are ms, s, m, h, d (24h) and
// w (7d); fractions are allowed (1.5h).

var (
	agoPattern = regexp.MustCompile(`^(?:\d+(?:\.\d+)?(?:ms|s|m|h|d|w))+$`)
	agoPart    = regexp.MustCompile(`(\d+(?:\.\d+)?)(ms|s|m|h|d|w)`)
	agoUnits   = map[string]time.Duration{
		"ms": time.Millisecond,
		"s":  time.Second,
		"m":  time.Minute,
		"h":  time.Hour,
		"d":  24 * time.Hour,
		"w":  7 * 24 * time.Hour,
	}
)

// parseTimeFlag reads an RFC 3339 timestamp, or a duration before now.
func parseTimeFlag(flag, value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	if ago, ok := parseAgo(value); ok {
		return time.Now().Add(-ago).UTC().Truncate(time.Second), nil
	}
	return time.Time{}, fmt.Errorf("invalid value %q for --%s: must be an RFC 3339 timestamp (e.g. 2025-02-27T15:04:05Z) or a duration ago (e.g. 24h, 7d)", value, flag)
}

// parseDateFlag reads a YYYY-MM-DD date, or a duration before now, which it
// turns into the UTC date that far back.
func parseDateFlag(flag, value string) (string, error) {
	if _, err := time.Parse(time.DateOnly, value); err == nil {
		return value, nil
	}
	if ago, ok := parseAgo(value); ok {
		return time.Now().Add(-ago).UTC().Format(time.DateOnly), nil
	}
	return "", fmt.Errorf("invalid value %q for --%s: must be YYYY-MM-DD (e.g. 2025-02-27) or a duration ago (e.g. 7d)", value, flag)
}

func parseAgo(value string) (time.Duration, bool) {
	if !agoPattern.MatchString(value) {
		return 0, false
	}
	var total time.Duration
	for _, part := range agoPart.FindAllStringSubmatch(value, -1) {
		n, err := strconv.ParseFloat(part[1], 64)
		if err != nil {
			return 0, false
		}
		total += time.Duration(n * float64(agoUnits[part[2]]))
	}
	return total, true
}
