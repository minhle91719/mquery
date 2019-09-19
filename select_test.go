package mquery

import (
	"testing"
)

func Test_selectQueryBuidler_ToQuery(t *testing.T) {
	qb := NewTable("user").Fields([]string{
		"id",
		"username",
		"password",
		"balance",
		"id_order_logging",
		"id_order",
		"status",
		"created_at",
		"updated_at",
		"deleted_at",
	})
	tests := []struct {
		name string
		sqb  IToQuery
		want string
	}{
		{
			name: "select with field",
			sqb:  qb.SelectBuilder().Fields("id", "username", "password"),
			want: "SELECT id,username,password FROM user",
		},
		{
			name: "select where and",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().Condition(NewCondition().And("id", Equal, 5))),
			want: "SELECT username,password FROM user WHERE id = 5",
		},
		{
			name: "select where or",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().Condition(NewCondition().And("id", Equal, 5))),
			want: "SELECT username,password FROM user WHERE id = 5",
		},
		{
			name: "select where and or mix",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().Condition(NewCondition().And("id", Equal, 5).And("password", Equal, "haha").Or("username", Equal, "hahaha"))),
			want: "SELECT username,password FROM user WHERE id = 5 AND password = \"haha\" OR username = \"hahaha\"",
		},
		{
			name: "select order by ASC",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().OrderByASC("username")),
			want: "SELECT username,password FROM user ORDER BY username ASC",
		}, {
			name: "select order by ASC",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().OrderByASC("username")),
			want: "SELECT username,password FROM user ORDER BY username ASC",
		},
		{
			name: "select JOIN",
			sqb: qb.SelectBuilder().NotCheckFieldValid().Fields("id_order_logging", "id_order", "status", "created_at").
				Where(qb.WhereBuilder().Condition(NewCondition().In(qb.SelectBuilder().Fields("DISTINCT(id_order)", "MAX(created_at)").
					Where(qb.WhereBuilder().GroupBy("id_order")).ToQuery(), "id_order", "created_at").And("status", Equal, "retry"))),
			want: "SELECT id_order_logging,id_order,status,created_at FROM user WHERE " +
				"(id_order,created_at) IN (SELECT distinct(id_order),max(created_at) FROM user GROUP BY id_order) AND status = \"retry\"",
		},
		{
			name: "select IN",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().Condition(NewCondition().In("?", "id")).Condition(NewCondition().Like("haha"))),
			want: "SELECT username,password FROM user WHERE (id) IN (1,2,3,4,5)",
		}, {
			name: "select IN Nested",
			sqb: qb.SelectBuilder().Fields("username", "password").Where(
				qb.WhereBuilder().Condition(NewCondition().In(qb.SelectBuilder().Fields("id").Where(qb.WhereBuilder().Condition(NewCondition().And("id", Equal, nil))).ToQuery(),
					"id"))),
			want: "SELECT username,password FROM user WHERE (id) IN (SELECT id FROM user WHERE id = ?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := tt.sqb
			if got := qb.ToQuery(); got != tt.want {
				t.Errorf("selectQueryBuidler.ToQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
