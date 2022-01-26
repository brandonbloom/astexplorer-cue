//go:build !js
// +build !js

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func main() {
	code, err := ioutil.ReadFile(os.Args[1])
	if err == nil {
		m := parseFile(string(code))
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(m)
	}
}
