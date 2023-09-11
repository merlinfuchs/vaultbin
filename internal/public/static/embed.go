package static

import "embed"

//go:embed *.js *.css *.ico
var FS embed.FS
