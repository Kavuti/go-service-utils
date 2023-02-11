package utils

import (
	"database/sql"
	"encoding/json"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Float64)
}
