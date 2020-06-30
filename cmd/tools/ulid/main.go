package main

import (
	"fmt"

	"bean/pkg/util"
)

func main() {
	id := util.Identifier{}

	fmt.Println("ULID: ", id.MustULID())
}
