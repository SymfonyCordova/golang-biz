package strings

import "strings"

func StringJoins(str1 string, str2 string) string {
	s := []string{str1, str2}
	return strings.Join(s, "")
}

func StringJoinSep(str1 string, str2 string, sep string) string {
	s := []string{str1, str2}
	return strings.Join(s, sep)
}
