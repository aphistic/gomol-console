package gomolconsole

import "github.com/aphistic/gomol"

/*
NewTemplateDefault will create a new default logging template.

Example output:
  [INFO] This is my message
*/
func NewTemplateDefault() *gomol.Template {
	tpl, _ := gomol.NewTemplate("[{{color}}{{ucase .LevelName}}{{reset}}] {{.Message}}")
	return tpl
}

/*
NewTemplateTimestamped will create a new logging template including the current timestamp.

Example output:
  2006-01-02 15:04:05.000 [INFO] This is my message
*/
func NewTemplateTimestamped() *gomol.Template {
	tpl, _ := gomol.NewTemplate("{{.Timestamp.Format \"2006-01-02 15:04:05.000\"}} [{{color}}{{ucase .LevelName}}{{reset}}] {{.Message}}")
	return tpl
}

/*
NewTemplateFull will create a new logging template including the current timestamp as well as any
attributes included with the message.

Example output:
  2006-01-02 15:04:05.000 [INFO] This is my message
     attr1: 1234
     attr2: value2
*/
func NewTemplateFull() *gomol.Template {
	tpl, _ := gomol.NewTemplate("{{.Timestamp.Format \"2006-01-02 15:04:05.000\"}} [{{color}}{{ucase .LevelName}}{{reset}}] {{.Message}}" +
		"{{if .Attrs}}{{range $key, $val := .Attrs}}\n   {{$key}}: {{$val}}{{end}}{{end}}")
	return tpl
}
