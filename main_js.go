package main

import (
	"encoding/json"
	"syscall/js"
)

func main() {
	block := make(chan struct{}, 0)
	js.Global().Set("__CUE_PARSE_FILE__", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		code := args[0].String()
		m := parseFile(code)
		res, _ := json.Marshal(m)
		return res
	}))
	<-block
}
