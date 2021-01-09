package main

import (
	"tools"
)

func main() {
	db := tools.GetDB()
	db.Close()
}
