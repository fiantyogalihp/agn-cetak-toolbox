package utils

import (
	"fmt"
	"strings"
)

// ToTerbilang converts a number to its spelled-out representation in Indonesian.
// The function takes an integer number and an optional suffix string.
// If a suffix is provided, it is appended to the spelled-out number.
// If two suffixes are provided, the first is used as the suffix, and the second can be "lower" or "upper" to control the capitalization.
//
// Example:
// ToTerbilang(int64(data.BillingAmount), "rupiah", "upper")
func ToTerbilang(num int64, suff ...string) string {
	switch len(suff) {
	case 1:
		return fmt.Sprintf("%s %s", hitTerbilang(num), suff[0])
	case 2:
		res := fmt.Sprintf("%s %s", hitTerbilang(num), suff[0])
		switch suff[1] {
		case "lower":
			return strings.ToLower(res)
		case "upper":
			return strings.ToUpper(res)
		}
		return res
	}

	return hitTerbilang(num)
}

// hitTerbilang converts an integer number to its spelled-out representation in Indonesian.
// It handles numbers up to 999,999,999,999,999 (quadrillion).
// The function recursively breaks down the number into thousands, millions, billions, and trillions,
// and combines the spelled-out representations of each part.
// The result is a string containing the spelled-out number.
func hitTerbilang(num int64) string {
	var s string
	satuan := [12]string{"", "satu", "dua", "tiga", "empat", "lima", "enam", "tujuh", "delapan", "sembilan", "sepuluh", "sebelas"}
	if num < 12 {
		s = satuan[num]
	} else if num < 20 {
		s = fmt.Sprintf("%s belas", hitTerbilang(num-10))
	} else if num < 100 {
		s = fmt.Sprintf("%s puluh %s", hitTerbilang(num/10), hitTerbilang(num%10))
	} else if num < 200 { // ratus
		s = fmt.Sprintf("seratus %s", hitTerbilang(num-100))
	} else if num < 1000 {
		s = fmt.Sprintf("%s ratus %s", hitTerbilang(num/100), hitTerbilang(num%100))
	} else if num < 2000 { // ribu
		s = fmt.Sprintf("seribu %s", hitTerbilang(num-1000))
	} else if num < 1000000 {
		s = fmt.Sprintf("%s ribu %s", hitTerbilang(num/1000), hitTerbilang(num%1000))
	} else if num < 2000000 { // juta
		s = fmt.Sprintf("satu juta %s", hitTerbilang(num-1000000))
	} else if num < 1000000000 {
		s = fmt.Sprintf("%s juta %s", hitTerbilang(num/1000000), hitTerbilang(num%1000000))
	} else if num < 2000000000 { // milyar
		s = fmt.Sprintf("satu milyar %s", hitTerbilang(num-1000000000))
	} else if num < 1000000000000 {
		s = fmt.Sprintf("%s milyar %s", hitTerbilang(num/1000000000), hitTerbilang(num%1000000000))
	} else if num < 2000000000000 { // triliun
		s = fmt.Sprintf("satu triliun %s", hitTerbilang(num-1000000000000))
	} else if num < 1000000000000000 {
		s = fmt.Sprintf("%s triliun %s", hitTerbilang(num/1000000000000), hitTerbilang(num%1000000000000))
	}
	return strings.TrimSpace(s)
}
