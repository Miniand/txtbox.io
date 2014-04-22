package main

import "github.com/Miniand/txtbox.io/web"

func main() {
	w, err := web.New()
	if err != nil {
		panic(err)
	}
	w.Run()
}
