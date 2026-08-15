package web

import (
	"embed"
	"io/fs"
)

//go:embed index.html
var files embed.FS

// FS отдаёт статику, вшитую в бинарник: деплой не зависит от рабочего каталога.
func FS() fs.FS { return files }
