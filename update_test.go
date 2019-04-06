package mquery

import "testing"

func Test_updateQueryBuilder_ToQuery(t *testing.T) {
	qb := NewTable("user").Fields(map[string]bool{
		"id":       true,
		"username": true,
		"password": true,
	})
	tests := []struct {
		name string
		uqb  IToQuery
		want string
	}{
		// TODO: Add test cases.
		{
			name: "update all",
			uqb: qb.UpdateBuilder().Value(map[string]interface{}{
				"username": "haha",
				"password": 5,
			}),
			want: "UPDATE user SET username = \"haha\",password = 5",
		},
		{
			name: "update where",
			uqb: qb.UpdateBuilder().Value(map[string]interface{}{
				"username": "haha",
				"password": 5,
			}).Where(qb.WhereBuilder().And("id", Equal, 5)),
			want: "UPDATE user SET username = \"haha\",password = 5 WHERE id = 5",
		}, {
			name: "update where nested",
			uqb: qb.UpdateBuilder().Value(map[string]interface{}{
				"username": "haha",
				"password": 5,
			}).Where(qb.WhereBuilder().And("id", Equal, 5).In("id", qb.SelectBuilder().Where(qb.WhereBuilder().And("id", Less, 5)))),
			want: `UPDATE user SET username = "haha",password = 5 WHERE id = 5 and id IN (SELECT * FROM user WHERE id < 5)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uqb.ToQuery(); got != tt.want {
				t.Errorf("updateQueryBuilder.ToQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
