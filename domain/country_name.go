package domain

import "fmt"

type countryName struct {
	common   string
	official string
}

func (cn *countryName) String() string {
	return fmt.Sprintf("countryName{common: '%s'; official: '%s'", cn.common, cn.official)
}
