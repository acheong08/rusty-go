package main

import (
	"fmt"

	"github.com/acheong08/rusty-go/genericutils"
)

func even(i int) (bool, error) {
	if (i % 2) == 0 {
		return true, nil
	}
	return false, fmt.Errorf("not even")
}

func evenWrapper(i int) (bool, error) {
	isEven, err := even(i)
	if err != nil {
		return genericutils.MakeGenericWithDefault[bool](), err
	}
	return isEven, nil
}

func main() {
	evenWrapper(1)
	evenWrapper(2)
}
