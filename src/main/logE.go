package main

import "fmt"

func logE(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format,a))
}
