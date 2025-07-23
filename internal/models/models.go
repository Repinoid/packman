package models

import "log/slog"

type sshConfStruct struct {
	Host       string `json:"host"`       // Адрес сервера
	User       string `json:"user"`       // Имя пользователя
	Password   string `json:"pasword"`    // Имя пользователя
	RemotePath string `json:"remotepath"` // Удалённый путь на сервере
}

var (
	Logger *slog.Logger
	SSHConf sshConfStruct
)

// Target для анмаршаллинга  {"path", "./archive_this2/*", "exclude": "*.tmp"}
type Target struct {
	Path    string `json:"path"`
	Exclude string `json:"exclude,omitempty"` // omitempty - необязательно есть маска
}

// Packet для  {"name": "packet-3", "ver": "<="2.0" } - куда паковать и условие для версии
type Packet struct {
	Name string `json:"name"`
	Ver  string `json:"ver"`
}

type Upack struct {
	Name    string   `json:"name"`
	Version string   `json:"ver"`     // будем считать что версия обязательна
	Targets []any    `json:"targets"` // any сиречь interface{}
	Packets []Packet `json:"packets"`
}

type Package struct {
	Name string `json:"name"`
	Ver  string `json:"ver,omitempty"` // omitempty - необязательно есть версия
}

/*
Пример файла пакета для упаковки:
packet.json

{
 "name": "packet-1",
 "ver": "1.10",
 "targets": [
  "./archive_this1/*.txt",
  {"path", "./archive_this2/*", "exclude": "*.tmp"},
 ]
 packets: {
  {"name": "packet-3", "ver": "<="2.0" },
 }
}

Пример файла для распаковки:

packages.json

{
 "packages": [
  {"name": "packet-1", "ver": ">=1.10"},
  {"name": "packet-2" },
  {"name": "packet-3", "ver": "<="1.10" },
 ]
}
*/
