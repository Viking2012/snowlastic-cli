package snowflake

import (
	"fmt"
	"regexp"
	"unicode"
)

func QuoteValue(i any) string {
	switch i.(type) {
	case string:
		return fmt.Sprintf("'%s'", i)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", i)
	case float32, float64:
		return fmt.Sprintf("%f", i)
	}
	return ""
}
func QuoteIdentifier(identifier string) string {
	var (
		needsQuotes                 bool
		prefix, _identifier, suffix string
		ret                         string
	)
	prefix, _identifier, suffix = extractIdentifier("", identifier, "")
	needsQuotes, _ = needsQuoting(_identifier)
	if needsQuotes {
		ret = fmt.Sprintf(`%s"%s"%s`, prefix, _identifier, suffix)
	} else {
		ret = fmt.Sprintf("%s%s%s", prefix, _identifier, suffix)
	}
	return ret
}
func extractIdentifier(prefix, identifier, suffix string) (string, string, string) {
	var ()
	re := regexp.MustCompile(`^(.+?\()(.+)(\)$)`)
	match := re.FindStringSubmatch(identifier)
	if len(match) > 0 {
		prefix, identifier, suffix = extractIdentifier(prefix+match[1], match[2], suffix+match[3])
	}
	return prefix, identifier, suffix
}

func needsQuoting(field string) (bool, error) {
	var (
		matched bool
		err     error
	)
	matched, err = regexp.Match(`^".+"$`, []byte(field))
	if err != nil {
		return true, err
	}
	if matched {
		return false, nil
	}
	matched, err = regexp.Match(`^[A-Za-z_].*`, []byte(field))
	if err != nil {
		return true, err
	}
	if !matched {
		return true, nil
	}

	// contains any non-alphanumeric character
	matched, err = regexp.Match(".*[^A-Za-z0-9_].*", []byte(field))
	if err != nil {
		return true, err
	}
	if matched {
		return true, nil
	}

	// Is not all uppercase
	if !isUpper(field) {
		return true, nil
	}

	return false, nil
}
func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
