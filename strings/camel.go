package strings

import "strings"

func LowerCamel(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToLower(str[:1]) + str[1:]
}

func ReceiverName(structName string) string {
	if len(structName) == 0 {
		return ""
	}
	return strings.ToLower(structName[:1])
}
