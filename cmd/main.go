package main

import "bean/pkg/infra"

func main() {
	container, err := infra.NewContainer("config.yaml")
	if nil != err {
		panic("failed creating container")
	}

	if err := container.ListenAndServe(); nil != err {
		panic(err)
	}
}
