package domain

import "fmt"

type Currency struct {
	Short  string
	Name   string
	Symbol string
}

func (c *Currency) String() string {
	return fmt.Sprintf("{'%s', '%s', '%s'}", c.Short, c.Name, c.Symbol)
}

type countryName struct {
	common   string
	official string
}

func (cn *countryName) String() string {
	return fmt.Sprintf("countryName{common: '%s'; official: '%s'", cn.common, cn.official)
}

type LatLng struct {
	Lat float32
	Lng float32
}
