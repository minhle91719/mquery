package mquery

import "testing"

func Test_updateQueryBuilder_ToQuery(t *testing.T) {
	qb := NewTable("user").Fields([]string{
		"id",
		"username",
		"password",
		"balance",
	})
	tests := []struct {
		name string
		uqb  IToQuery
		want string
	}{
		// TODO: Add test cases.
		{
			name: "update all",
			uqb:  qb.UpdateBuilder().Fields("username", "password"),
			want: "UPDATE user SET username = ?,password = ?",
		},
		{
			name: "update where",
			uqb: qb.UpdateBuilder().Values(map[string]interface{}{
				"username": "haha",
				"password": 5,
			}).Where(qb.WhereBuilder().Condition(NewCondition().And("id", Equal, 5))),
			want: "UPDATE user SET username = \"haha\",password = 5 WHERE id = 5",
		}, {
			name: "update where nested",
			want: `UPDATE user SET username = "haha",password = 5 WHERE id = 5 AND (id) IN (SELECT id FROM user WHERE id < 5)`,
			uqb: qb.UpdateBuilder().Values(map[string]interface{}{
				"username": "haha",
				"password": 5,
			}).Where(qb.WhereBuilder().Condition(NewCondition().And("id", Equal, 5).
				In(qb.SelectBuilder().Fields("id").Where(qb.WhereBuilder().Condition(NewCondition().And("id", Less, 5))).ToQuery(), "id"))),
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
