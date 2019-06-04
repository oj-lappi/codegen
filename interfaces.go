package bakefiles

import (
	"html/template"
)

type Parser interface {
	parse(string) interface{}
	BakeTemplate(template.Template, string) []byte
	BakeFile(string, string) []byte //From filenames
}
