package mquery

import (
	"testing"
)

func Test_selectQueryBuidler_ToQuery(t *testing.T) {
	qb := NewQueryBuilder("user", Column(
		"id",
		"username",
		"password",
		"balance",
		"id_order_logging",
		"id_order",
		"status",
		"created_at",
		"updated_at",
		"deleted_at"))
	
	tests := []struct {
		name string
		sqb  toQuery
		want string
	}{
		{
			name: "select with field",
			sqb:  qb.Select(SelectField("id","username","password")).Where(Condition(AndPair([]interface{}{
				"status","status","status",
			},1,2,3))),
			want: "SELECT id,username,password FROM user",
		},
		{
			name: "select where and",
			sqb:  qb.Select(SelectField("username", "password")).Where(Condition(And("id", EqualOps, 5))),
			want: "SELECT username,password FROM user WHERE id = 5",
		},
		{
			name: "select where or",
			sqb:  qb.Select(SelectField("username", "password")).Where(Condition(Or("id", EqualOps, 5))),
			want: "SELECT username,password FROM user WHERE id = 5",
		},
		{
			name: "select where and or mix",
			sqb:  qb.Select(SelectField("username", "password")).Where(Condition(And("id", EqualOps, 5), And("password", EqualOps, "haha"), Or("username", EqualOps, "hahaha"))),
			want: "SELECT username,password FROM user WHERE id = 5 AND password = \"haha\" OR username = \"hahaha\"",
		},
		{
			name: "select order by ASC",
			sqb:  qb.Select(SelectField("username", "password")).Where(OrderBy("username", ASC)),
			want: "SELECT username,password FROM user ORDER BY username ASC",
		}, {
			name: "select order by DESC",
			sqb:  qb.Select(SelectField("username", "password")).Where(OrderBy("username", DESC)),
			want: "SELECT username,password FROM user ORDER BY username DESC",
		},
		{
			name: "select IN 2 key",
			sqb: qb.Select(SelectField("id_order_logging", "id_order", "status", "created_at")).
				Where(
					Condition(
						In(qb.Select(SelectField(Distinct("id_order"), Max("created_at"))).Where(GroupBy("id_order")).ToQuery(), "id_order", "created_at"),
						And("status", EqualOps, "retry")),
				),
			want: "SELECT id_order_logging,id_order,status,created_at FROM user WHERE " +
				"(id_order,created_at) IN (SELECT DISTINCT(id_order),MAX(created_at) FROM user GROUP BY id_order) AND status = \"retry\"",
		},
		{
			name: "select IN",
			sqb:  qb.Select(SelectField("username", "password")).Where(Condition(In("1,2,3,4,5", "id"))),
			want: "SELECT username,password FROM user WHERE (id) IN (1,2,3,4,5)",
		}, {
			name: "select IN Nested",
			sqb: qb.Select(SelectField("username", "password")).
				Where(Condition(In(qb.Select(SelectField("id")).Where(Condition(And("id", EqualOps))).ToQuery(), "id"))),
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
