В папке cmd/packman main.go
Перейти в эту папку
***cd cmd/packman***

Поднять SSL server
***docker compose up -d***

Далее как в техзадании 
```
pm create ./packet.json

pm update ./packages.json
```

конфигурации в packet.json и packages.json

Если принципиально, могу добавить тесты, с Suite & testcontainer


*Сделать на GO пакетный менеджер

должен уметь упаковывать файлы в архив, и заливать их на сервер по SSH
должен уметь скачивать файлы архивов по SSH и распаковывать.

Фаил для упаковки должен иметь формат .yaml или json
в файле должны быть указаны пути по которым нужно подобрать файлы по маске

Пример файла пакета для упаковки:*

packet.json
```
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
```
Пример файла для распаковки:

packages.json
```
{
 "packages": [
  {"name": "packet-1", "ver": ">=1.10"},
  {"name": "packet-2" },
  {"name": "packet-3", "ver": "<="1.10" },
 ]
}
```
*Сделать commandline tools с командами:

pm create ./packet.json

pm update ./packages.json

PS: Можно использовать любые допущения которые сделают разработку тестового задания проще*

