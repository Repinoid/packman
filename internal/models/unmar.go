package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

// свой анмаршаллинг типа Upack
func (upa *Upack) UnmarshalJSON(data []byte) (err error) {

	type alias Upack
	aux := &struct {
		*alias
	}{
		alias: (*alias)(upa),
	}

	err = json.Unmarshal([]byte(data), aux)

	// Должны присутствовать все поля
	if upa.Name == "" {
		return errors.New("no \"Name\" field")
	}
	if upa.Version == "" {
		return errors.New("no \"Version\" field")
	}
	if len(upa.Targets) == 0 {
		return errors.New("no \"Targets\" field")
	}
	if len(upa.Packets) == 0 {
		return errors.New("no \"Packets\" field")
	}

	for i, t := range upa.Targets {
		switch v := t.(type) {

		case string:
			// из строки - в map[string] с ключом "path" без "exclude"
			upa.Targets[i] = map[string]string{"path": v}

		case map[string]any:
			_, ok := v["path"]
			if !ok {
				// если нет ключа "path" - возврат по ошибке
				return fmt.Errorf("wrong key in Targets %v", v)
			}
			_, ok = v["exclude"]
			// transform to map[string]string
			if !ok {
				// если нет ключа "exclude"
				upa.Targets[i] = map[string]string{"path": v["path"].(string)}
			} else {
				upa.Targets[i] = map[string]string{"path": v["path"].(string), "exclude": v["exclude"].(string)}
			}

		}
	}
	return
}
