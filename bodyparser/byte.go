package bodyparser

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var byteMatch = regexp.MustCompile(`(?i)^\+?(\d*)(\.\d*)? *(b|kb|mb|gb|tb|pb)$`)
var byteUnitMap = map[string]int64{
	"b":  1,
	"kb": 1 << 10,
	"mb": 1 << 20,
	"gb": 1 << 30,
	"tb": 1 << 40,
	"pb": 1 << 50,
}

func parseByte(limit any) int64 {
	if l, ok := limit.(int64); ok {
		return l
	} else if l, ok := limit.(int); ok {
		return int64(l)
	} else if l, ok := limit.(string); ok {
		matches := byteMatch.FindStringSubmatch(l)
		var num float64 = 0
		var unit string = "b"
		var err error = nil

		if len(matches) == 4 {
			num, err = strconv.ParseFloat(matches[1]+matches[2], 64)
			unit = strings.ToLower(matches[3])
		} else {
			num, err = strconv.ParseFloat(l, 64)
		}

		if err != nil {
			panic(errors.New("limit is malformed"))
		}

		result := num * float64(byteUnitMap[unit])

		if result >= math.MaxInt64 || result <= math.MinInt64 {
			panic(errors.New("limit is too huge"))
		}

		return int64(result)
	}

	panic(errors.New("limit is malformed"))
}
