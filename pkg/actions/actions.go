package actions

import "fmt"

func SetOutput(name string, value interface{}) {
	fmt.Printf("::set-output name=%s::%v\n", name, value)
}
