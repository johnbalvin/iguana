package files

import (
	"bytes"
	"regexp"
)

var regexCorrectCSs = regexp.MustCompile(`url\([\w\.\-_/ ]+?((\.js)|(\.json)|(\.woff)|(\.woff2)|(\.ttf)|(\.json)|(\.mjs)|(\.txt)|(\.css)|(\.md)|(\.gif)|(\.ico)|(\.jpeg)|(\.jpg)|(\.png)|(\.tiff)|(\.tif)|(\.svg))\)`)

//before, the script didn't catch this format: url(../img/icons/ic-search.svg)
//so I will just convert it to url("../img/icons/ic-search.svg")
func CorrectFormatCss(content []byte) []byte {
	return regexCorrectCSs.ReplaceAllFunc(content, func(current []byte) []byte {
		current = bytes.ReplaceAll(current, []byte("url("), []byte(`url("`))
		current = bytes.ReplaceAll(current, []byte(")"), []byte(`")`))
		return current
	})
}
