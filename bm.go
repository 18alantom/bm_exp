package main

import "fmt"

type BMConfig struct {
	Cache string
}

func main() {
	config := BMConfig{"/Users/alan/Desktop/code/test_go/bm_poc/.cache"}
	_ = config

	fmt.Println("Hello, World")
}
