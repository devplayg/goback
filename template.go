package goback

import (
	"fmt"
	"io/ioutil"
)

func DisplayWithLocalFile(name string) string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tpl/%s.html", name))
	if err != nil {
		log.Error(err)
	}
	return string(b)
}
