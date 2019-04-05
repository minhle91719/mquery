package mquery

import "testing"

func Test_insertQueryBuilder_ToQuery(t *testing.T) {
	type fields struct {
		tableName string
		col       []string
		value     []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// iqb := &insertQueryBuilder{
			// 	tableName: tt.fields.tableName,
			// 	col:       tt.fields.col,
			// 	value:     tt.fields.value,
			// }
			// if got := iqb.ToQuery(); got != tt.want {
			// 	t.Errorf("insertQueryBuilder.ToQuery() = %v, want %v", got, tt.want)
			// }
		})
	}
}
