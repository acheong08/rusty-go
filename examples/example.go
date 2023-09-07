package main

import (
  "fmt"
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
		return makeGenericWithDefault[bool](), err
	}
  return isEven, nil
}

func main() {
  evenWrapper(1)
  evenWrapper(2)
}
func makeGenericWithDefault[T any]() T {
				var t T
				return t
}