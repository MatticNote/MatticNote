package mn_template

import "embed"

//go:embed *.django */*.django
var Templates embed.FS
