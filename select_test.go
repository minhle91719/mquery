package mquery

import (
	"testing"
)

func Test_selectQueryBuidler_ToQuery(t *testing.T) {
	qb := NewTable("user").Fields(map[string]bool{
		"id":       true,
		"username": true,
		"password": true,
	})
	tests := []struct {
		name string
		sqb  IToQuery
		want string
	}{
		{
			name: "select all",
			sqb:  qb.SelectBuilder(),
			want: "SELECT * FROM user",
		},
		{
			name: "select with field",
			sqb:  qb.SelectBuilder().Fields("id", "username", "password"),
			want: "SELECT id,username,password FROM user",
		},
		{
			name: "select where and",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().And("id", Equal, 5)),
			want: "SELECT username,password FROM user WHERE id = 5",
		},
		{
			name: "select where or",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().Or("id", Equal, 5)),
			want: "SELECT username,password FROM user WHERE id = 5",
		},
		{
			name: "select where and or mix",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().And("id", Equal, 5).Or("username", Equal, "hahaha").And("password", Equal, "haha")),
			want: "SELECT username,password FROM user WHERE id = 5 and password = \"haha\" or username = \"hahaha\"",
		},
		{
			name: "select order by ASC",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().OrderByASC("username")),
			want: "SELECT username,password FROM user ORDER BY username ASC",
		}, {
			name: "select order by ASC",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().OrderByASC("username")),
			want: "SELECT username,password FROM user ORDER BY username ASC",
		}, {
			name: "select JOIN",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().Join("account", "id_user", "id_user").And("balance", Equal, 5)),
			want: "SELECT username,password FROM user JOIN account ON user.id_user = account.id_user WHERE balance = 5",
		}, {
			name: "select IN",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().In("id_user", 1, 2, 3, 4, 5)),
			want: "SELECT username,password FROM user WHERE id_user IN (1,2,3,4,5)",
		}, {
			name: "select IN Nested",
			sqb:  qb.SelectBuilder().Fields("username", "password").Where(qb.WhereBuilder().In("id_user", qb.SelectBuilder().Fields("id").Where(qb.WhereBuilder().And("id", Equal, 5)))),
			want: "SELECT username,password FROM user WHERE id_user IN (SELECT id FROM user WHERE id = 5)",
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
