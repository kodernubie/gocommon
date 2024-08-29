package main

import (
	"fmt"

	"github.com/kodernubie/gocommon/conf"
)

func main() {

	fmt.Println("Test Config")
	fmt.Println("String :", conf.Str("CONFIG_STRING"))
	fmt.Println("String with Default:", conf.Str("CONFIG_STRING_DEFAULT", "Not provided"))

	fmt.Println("Int :", conf.Int("CONFIG_INT"))
	fmt.Println("Int with Default:", conf.Int("CONFIG_INT_DEFAULT", -1))

	fmt.Println("Float :", conf.Float("CONFIG_FLOAT"))
	fmt.Println("Float with Default:", conf.Float("CONFIG_FLOAT_DEFAULT", -99))

	fmt.Println("Bool :", conf.Bool("CONFIG_BOOL"))
	fmt.Println("Bool with Default:", conf.Bool("CONFIG_STRING_DEFAULT", false))
}
