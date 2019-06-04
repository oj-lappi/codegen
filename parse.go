package bakefiles

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"path/filepath"
	"reflect"
)

type BulkString struct {
	Files    string `toml:"files,omitempty"`
	Literals string `toml:"literals,omitempty"`
}

type BulkStringConfig map[string]BulkString

func parseConfig(filePath string, v interface{}) {

	//Read into a temporary format
	var temp_config BulkStringConfig

	if _, err := toml.DecodeFile(filePath, &temp_config); err != nil {
		panic(err.Error())
	}
	fmt.Println("config:", temp_config)

	//TODO:ensure v is a pointer, otherwise the value won't be settable
	rv := reflect.ValueOf(v)

	//Get all bulkstrings merged and shove them into the pointers fields

	//toml package uses indirect, do we need that too?
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	} else {
		fmt.Println("Value passed to parseConfig is not a pointer")
		return
		//We have a problem (I think)
	}
	//We have two cases (I assume)
	//structs, and maps

	switch rv.Kind() {
	case reflect.Ptr:
		//Need recursion for this, which means adding a layer
		fmt.Println("Pointers not supported")
	case reflect.Struct:
		s := rv
		t := s.Type()
		for i := 0; i < s.NumField(); i++ {
			fieldName := t.Field(i).Name
			configVal, ok := temp_config[fieldName]
			if !ok {
				fmt.Printf("No field called %s in toml file\n", fieldName)
				continue
			}

			//f := s.Field(i)
			//fieldType := f.Type()

			//Get name value etc

			//Get files, and merge them with literals
			fileStrings := readFiles(configVal.Files)
			literalStrings := configVal.Literals
			strings := append(fileStrings, literalStrings)
			fmt.Printf("all %s:\n", fieldName)
			fmt.Printf("%v\n", strings)
			//Assign value to field
			//TODO: implement

			//TODO:GET THIS WORKING WITH TEMPLATES ASAP
		}
	case reflect.Map:
		fmt.Println("maps not supported")
	default:
		fmt.Println(rv.Type(), "not supported")
	}

}

//returns file contents
func readFiles(filePath string) []string {
	fileNames, err := filepath.Glob(filePath)
	if err != nil {
		panic(err.Error())
	}
	files := make([]string, 2)
	for _, fileName := range fileNames {
		file, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err.Error())
		}
		files = append(files, string(file))
	}
	return files
	//TODO: get fancy and make a map[string][string] instead
	//Or possibly structs, but is that too coupled?
}

/*
Goal:

Say static/assets/svg includes the following files:

	a.svg
	b.svg
	c.svg

bake.toml includes:

	[svg]
		files 	= static/assets/svg/*
		literal = ['''
			<svg>
			...
			</svg>''']

And bake is a struct:


	type bake struct{
		svg []string
	}

Then calling parseConfig("bake.toml",&bake{}) will produce:

	{
		"<svg>... </svg>",//a.svg
		"<svg>... </svg>",//b.svg
		"<svg>... </svg>",//c.svg
		"<svg>... </svg>"//The literal string
	}

*/
