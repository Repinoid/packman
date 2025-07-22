package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_unmar(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		isErr  bool
		errStr string
	}{
		{
			// "target" вместо "targets"
			name: "Кривая структура",
			data: `{
				"name": "packet-1",
				"ver": "1.10",
				"target": [
					"./archive_this1/*.txt",
					{"path": "./archive_this2/*", "exclude": "*.tmp"}
						],
				"packets": [
				{"name": "packet-3", "ver": "<=2.0"}
				]
				}`,
			isErr:  true,
			errStr: "no \"Targets\" field",
		},
		{
			name: "Правильная структура",
			data: `{
				"name": "packet-1",
				"ver": "1.10",
				"targets": [
					"./archive_this1/*.txt",
					{"path": "./archive_this2/*", "exclude": "*.tmp"}
						],
				"packets": [
				{"name": "packet-3", "ver": "<=2.0"}
				]
				}`,
			isErr: false,
		},
		{
			name: "Нет ключа path ",
			data: `{
				"name": "packet-1",
				"ver": "1.10",
				"targets": [
					"./archive_this1/*.txt",
					{"patha": "./archive_this2/*", "exclude": "*.tmp"}
						],
				"packets": [
				{"name": "packet-3", "ver": "<=2.0"}
				]
				}`,
			isErr:  true,
			errStr: "wrong key in Targets",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := unmar([]byte(tt.data))
			if tt.isErr {
				assert.ErrorContains(t, err, tt.errStr)
			} else {
				assert.NilError(t, err, "no error expected, man")
			}

		})
	}
}
