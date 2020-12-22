
package main

import (
	"stockplay/core"
)
const (
	download = false
)
func main() {
	sp := core.NewStockPlay()
	if download {
		err := sp.Download()
		if err != nil {
			panic(err)
		}
	}
	err := sp.Start()
	if err != nil {
		panic(err)
	}
}