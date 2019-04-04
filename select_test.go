package mquery

import (
	"fmt"
	"testing"
)

func Test_selectQueryBuidler_ToQuery(t *testing.T) {

	tests := []struct {
		name string
		sqb  SelectQueryBuilder
		want string
	}{
		{
			name: "select all",
			sqb:  Select("user"),
			want: "SELECT * FROM user",
		},
		{
			name: "select with field",
			sqb:  Select("user").Fields("id", "username", "password"),
			want: "SELECT id,username,password FROM user",
		},
		{
			name: "select where and",
			sqb:  Select("user").Fields("username", "password").And("id", 5),
			want: "SELECT username,password FROM user WHERE id = 5",
		},
		{
			name: "select where or",
			sqb:  Select("user").Fields("username", "password").Or("id", 5),
			want: "SELECT username,password FROM user WHERE id = 5",
		},
		{
			name: "select where and or mix",
			sqb:  Select("user").Fields("username", "password").And("id", 5).Or("username", "hahaha").And("password", "haha"),
			want: "SELECT username,password FROM user WHERE id = 5 and password = \"haha\" or username = \"hahaha\"",
		},
		{
			name: "select order by ASC",
			sqb:  Select("user").Fields("username", "password").OrderByASC("username"),
			want: "SELECT username,password FROM user ORDER BY username ASC",
		}, {
			name: "select order by ASC",
			sqb:  Select("user").Fields("username", "password").OrderByASC("username"),
			want: "SELECT username,password FROM user ORDER BY username ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := tt.sqb
			if got := qb.ToQuery(); got != tt.want {
				t.Errorf("selectQueryBuidler.ToQuery() = %v, want %v", got, tt.want)
			} else {
				fmt.Println(got)
			}
		})
	}
}
