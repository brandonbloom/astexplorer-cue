package main

import (
	"fmt"
	"os"
	"reflect"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
)

func parseFile(code string) map[string]interface{} {
	f, _ := parser.ParseFile("", code)
	return walk(f)
}

func walk(node interface{}) map[string]interface{} {
	if node == nil {
		return nil
	}

	m := make(map[string]interface{})

	val := reflect.ValueOf(node)
	if val.IsNil() {
		return nil
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	ty := val.Type()
	m["_type"] = ty.Name()
	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		val := val.Field(i)
		// field.IsExported not available in Go 1.13.
		if field.PkgPath == "" {
			continue
		}
		if field.Type.Name() == "Pos" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.Array, reflect.Slice:
			list := make([]interface{}, 0, val.Len())
			for i := 0; i < val.Len(); i++ {
				if item := walk(val.Index(i).Interface()); item != nil {
					list = append(list, item)
				}
			}
			m[field.Name] = list
		case reflect.Ptr:
			if child := walk(val.Interface()); child != nil {
				m[field.Name] = child
			}
		case reflect.Interface:
			if child := walk(val.Interface()); child != nil {
				m[field.Name] = child
			}
		case reflect.String:
			m[field.Name] = val.String()
		case reflect.Int:
			if field.Type.Name() == "Token" {
				m[field.Name] = token.Token(val.Int()).String()
			} else {
				m[field.Name] = val.Int()
			}
		case reflect.Bool:
			m[field.Name] = val.Bool()
		default:
			fmt.Fprintln(os.Stderr, field)
		}
	}
	if n, ok := node.(ast.Node); ok {
		start := n.Pos().Offset()
		end := n.End().Offset()
		m["Loc"] = map[string]interface{}{"Start": start, "End": end}
	}
	return m
}
