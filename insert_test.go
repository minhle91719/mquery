package mquery

import (
	"fmt"
	"testing"
)

func Test_insertQueryBuilder_ToQuery(t *testing.T) {
	test := NewQueryBuilder("haha", Column("perm"))
	var condition = []ConditionOption{}
	condition = append(condition, Pair(false, []interface{}{
		"perm",
	}, "read"))
	fmt.Println(test.Update(NotCheckFieldUpdate(),UpdateField("perm", NULL())).Where(Condition(condition...)).ToQuery())
	fmt.Println()
	//test.Select()
	qb := NewQueryBuilder("user", Column("id",
		"username",
		"password",
		"balance"))
	tests := []struct {
		name string
		iqb  toQuery
		want string
	}{
		{
			name: "insert value",
			iqb: qb.Insert(WithField("username", "password"), OnDuplicate(map[interface{}]interface{}{
				"username": "hahah",
			})),
			want: `INSERT INTO user(username,password) VALUE(?,?)`,
		},
		{
			name: "insert 2 value",
			iqb:  qb.Insert(WithField("username", "password"), WithValues(2)),
			want: `INSERT INTO user(username,password) VALUES(?,?),(?,?)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.iqb.ToQuery(); got != tt.want {
				t.Errorf("insertQueryBuilder.ToQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
