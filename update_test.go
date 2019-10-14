package mquery

import "testing"

func Test_updateQueryBuilder_ToQuery(t *testing.T) {
	qb := NewQueryBuilder("user", Column(
		"id",
		"username",
		"password",
		"balance"))
	tests := []struct {
		name string
		uqb  toQuery
		want string
	}{
		// TODO: Add test cases.
		{
			name: "update all",
			uqb:  qb.Update(UpdateField("username", Now()), UpdateField("password", nil)).Where(),
			want: "UPDATE user SET username = ?,password = ?",
		},
		{
			name: "update where",
			uqb:  qb.Update(UpdateField("username", "haha"), UpdateField("password", 5)).Where(Condition(And("id", EqualOps, 5))),
			want: "UPDATE user SET username = \"haha\",password = 5 WHERE id = 5",
		}, {
			name: "update where nested",
			want: `UPDATE user SET username = "haha",password = 5 WHERE id = 5 AND (id) IN (SELECT id FROM user WHERE id < 5)`,
			uqb: qb.Update(UpdateField("username", "haha"), UpdateField("password", 5)).Where(
				//Condition(And("id", EqualOps, 5), In(qb.Select(SelectField("id")).From().Where(Condition(And("id", LessOps, 5))).ToQuery(), "id")),
			),
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
