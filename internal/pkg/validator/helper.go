package validator

import "regexp"

var (
	EmailRX           = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	EgyPhoneNumbersRX = regexp.MustCompile("^(\\+20|0)1[0125]\\d{8}$")
)

func Matches(word string, regex regexp.Regexp) bool {
	return regex.MatchString(word)
}

func In[T comparable](item T, items ...T) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}

	return false
}
