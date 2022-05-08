package check

import "strconv"

//CheckDigits is a function to check digits
func CheckDigits(s string) bool {
	n, _ := strconv.Atoi(s)
	count := 0
	for n > 0 {
		n = n / 10
		count++
	}
	if count == 10 {
		return true
	}
	return false
}
