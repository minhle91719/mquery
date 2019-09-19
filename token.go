package mquery

import (
	"fmt"
	"strings"
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
func Count(col interface{}) string {
	return fmt.Sprintf("COUNT(%s)", col)
}
