package platform

import (
	"fmt"
	"math"
	"strconv"
)

func ByteCountDecimal(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

func ByteFormat(inputNum float64, precision int) string {

	if precision <= 0 {
		precision = 4
	}

	var unit string
	var returnVal float64

	if inputNum >= 1000000000000000000000000 {
		returnVal = RoundUp((inputNum / 1208925819614629174706176), precision)
		unit = " YB" // yottabyte
	} else if inputNum >= 1000000000000000000000 {
		returnVal = RoundUp((inputNum / 1180591620717411303424), precision)
		unit = " ZB" // zettabyte
	} else if inputNum >= 10000000000000000000 {
		returnVal = RoundUp((inputNum / 1152921504606846976), precision)
		unit = " EB" // exabyte
	} else if inputNum >= 1000000000000000 {
		returnVal = RoundUp((inputNum / 1125899906842624), precision)
		unit = " PB" // petabyte
	} else if inputNum >= 1000000000000 {
		returnVal = RoundUp((inputNum / 1099511627776), precision)
		unit = " TB" // terrabyte
	} else if inputNum >= 1000000000 {
		returnVal = RoundUp((inputNum / 1073741824), precision)
		unit = " GB" // gigabyte
	} else if inputNum >= 1000000 {
		returnVal = RoundUp((inputNum / 1048576), precision)
		unit = " MB" // megabyte
	} else if inputNum >= 1000 {
		returnVal = RoundUp((inputNum / 1024), precision)
		unit = " KB" // kilobyte
	} else {
		returnVal = inputNum
		unit = " bytes" // byte
	}

	return strconv.FormatFloat(returnVal, 'f', precision, 64) + unit

}
