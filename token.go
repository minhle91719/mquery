package mquery

import (
	"fmt"
	"html"
	"strings"
	"time"
)

var replaceToken = strings.NewReplacer(
	"DISTINCT", "",
	"MAX", "", "(", "", ")", "",
	"COUNT", "")

func Distinct(col interface{}) string {
	return fmt.Sprintf("DISTINCT(%s)", col)
}
func Max(col interface{}) string {
	return fmt.Sprintf("MAX(%s)", col)
}
func Now() string {
	return "NOW()"
}
func NULL() string {
	return "NULL"
}

var listToken = []string{
	"NOW()", "NULL",
}

func interfaceToString(value interface{}) string {
	result := ""
	switch value.(type) {
	case int, uint:
		result = fmt.Sprintf("%d", value)
	case string:
		values := fmt.Sprintf("%s", value)
		for _, v := range listToken {
			if v == strings.ToUpper(values) {
				return v
			}
		}
		result = fmt.Sprintf(`"%s"`, html.EscapeString(fmt.Sprintf("%s", value)))
	case time.Time:
		result = `"` + value.(time.Time).String() + `"`
	case bool:
		result = fmt.Sprint(value)
	case nil:
		result = "?"
	default:
		return fmt.Sprint(value)
	}
	return result
}
func genValueParam(length int) (value string) {
	listValue := make([]string, 0, length)
	for i := 0; i < length; i++ {
		listValue = append(listValue, "?")
	}
	return "(" + strings.Join(listValue, ",") + ")"
}
