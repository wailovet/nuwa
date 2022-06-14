package main

import (
	"fmt"

	"github.com/wailovet/nuwa"
)

type Test struct {
	Id  int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Key string `json:"key"`
}

func main() {
	nuwa.Sqlited().Config("./test.db", "nuwa_")
	fmt.Println(nuwa.Sqlited().Xorm().Sync2(&Test{}))
}
