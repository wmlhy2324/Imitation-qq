package ctype

import (
	"database/sql/driver"
	"encoding/json"
)

type SystemMsg struct {
	Type uint8 `json:"type"` //1涉黄 2涉恐 3涉政 4不当言论

}

func (c *SystemMsg) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}
func (c SystemMsg) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
