package main

import (
	"fmt"

	"bean/components/scalar"
)

func main() {
	id := scalar.Identifier{}

	fmt.Println("ULID: ", id.ULID())
}
