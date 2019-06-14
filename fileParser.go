package main

import (
	"io/ioutil"
)

//ReadTextFile reads a given text file
func ReadTextFile(path string) string {

	text, err := ioutil.ReadFile(path)
	CheckError(err)

	return string(text)
}
