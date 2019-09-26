package mquery

import (
	"fmt"
	"strings"
)

type Condition interface {
	And(key interface{}, ops Operator, value interface{}) Condition
	Or(key interface{}, ops Operator, value interface{}) Condition
	Like(key string) Condition
	
	In(queryIn string, colIn ...string) Condition
	
	ToCondititonString() string
}

type conditionImpl struct {
	and   []string
	or    []string
	param int
}

func (c *conditionImpl) Like(key string) Condition {
	c.and = append(c.and, fmt.Sprintf("%s %s ?", key, Like))
	return c
}

func (c *conditionImpl) And(key interface{}, ops Operator, value interface{}) Condition {
	con := ""
	if value == nil {
		con = fmt.Sprintf("%s %s ?", key, ops)
	} else {
		con = fmt.Sprintf("%s %s %s", key, ops, interfaceToString(value))
	}
	c.and = append(c.and, con)
	return c
}

func (c *conditionImpl) Or(key interface{}, ops Operator, value interface{}) Condition {
	con := ""
	if value == nil {
		con = fmt.Sprintf("%s %s ?", ops, key)
	} else {
		con = fmt.Sprintf("%s %s %s", key, ops, interfaceToString(value))
	}
	c.or = append(c.or, con)
	return c
}

func (c *conditionImpl) In(queryIn string, colIn ...string) Condition {
	con := ""
	colValue := "(" + strings.Join(colIn, ",") + ")"
	if queryIn == "" {
		con = fmt.Sprintf("%s IN (%s)", colValue, "?")
	} else {
		con = fmt.Sprintf("%s IN (%s)", colValue, queryIn)
	}
	c.and = append(c.and, con)
	return c
}

func (c *conditionImpl) ToCondititonString() string {
	con := []string{}
	if c.and != nil {
		con = append(con, strings.Join(c.and, " AND "))
	}
	if c.or != nil {
		con = append(con, strings.Join(c.or, " OR "))
	}
	if con == nil {
		return ""
	}
	return strings.Join(con, " OR ")
}

func NewCondition() Condition {
	return &conditionImpl{}
}
