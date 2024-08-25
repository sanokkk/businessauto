package uuid_helpers

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

// Value - метод для преобразования UUIDArray в JSON при сохранении в базу данных
func (u UUIDArray) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Scan - метод для чтения JSON из базы данных и преобразования его в UUIDArray
func (u *UUIDArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert %v to UUIDArray", value)
	}

	return json.Unmarshal(bytes, u)
}
