package main

import (
	"fmt"

	"github.com/denisbrodbeck/machineid"
)

func main() {
	id, err := machineid.ProtectedID("**TaskMaster**")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(id)
}
