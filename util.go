package peajob

import "github.com/teris-io/shortid"

//GenShortID 获取随机id
func GenShortID() (string, error) {
	return shortid.Generate()
}
