package helpers

import "github.com/google/uuid"

func ConvertStringArray(uuids []string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(uuids))

	for _, id := range uuids {
		if id == "" {
			continue
		}
		parsed, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		result = append(result, parsed)
	}

	return result, nil
}
