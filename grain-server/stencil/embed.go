package stencil

import "embed"

//go:embed web/*.grain
var WebTemplateFS embed.FS

//go:embed mysql/*.grain
var MysqlTemplateFS embed.FS

//go:embed mysql/admin/*.grain
var MysqlAdminTemplateFS embed.FS

//go:embed flutter/*.grain
var FlutterTemplateFS embed.FS
