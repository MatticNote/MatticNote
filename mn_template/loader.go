package mn_template

import "embed"

//go:embed *.pug */*.pug
var Templates embed.FS
