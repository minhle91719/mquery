package mquery

import "testing"

func Test_insertQueryBuilder_ToQuery(t *testing.T) {
	qb := NewTable("user").Fields(map[string]bool{
		"id":       false,
		"username": true,
		"password": true,
		"boolss":   true,
	})
	tests := []struct {
		name string
		iqb  IToQuery
		want string
	}{
		{
			name: "insert",
			iqb:  qb.InsertBuilder().Value("username", "password"),
			want: `INSERT INTO user(username,password) VALUES("minhle","deptrai")`,
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
