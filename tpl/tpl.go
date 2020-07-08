package tpl

import (
	"io/ioutil"
)

func displayWithLocalFile(path string) string {
	b, _ := ioutil.ReadFile(path)
	return string(b)
}
