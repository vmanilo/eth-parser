package internal

import (
	"regexp"
	"strconv"
)

func hexToInt(hex string) int64 {
	val, _ := strconv.ParseInt(hex[2:], 16, 64)
	return val
}

func intToHex(val int64) string {
	return "0x" + strconv.FormatInt(val, 16)
}

var hexRegexp = regexp.MustCompile("^0x[a-f0-9]+$")

func isValidHex(hex string) bool {
	return hexRegexp.MatchString(hex)
}
