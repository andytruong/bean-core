package main

import (
	"fmt"

	"bean/components/unique"
)

func main() {
	id := unique.Identifier{}

	fmt.Println("ULID: ", id.MustULID())
}
