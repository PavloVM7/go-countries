package domain

import "strings"

type countryExt struct {
	borders []string
}

func (c *countryExt) Borders() []string {
	if len(c.borders) == 0 {
		return []string{}
	}
	result := make([]string, len(c.borders))
	copy(result, c.borders)
	return result
}
func (c *countryExt) SetBorders(borders ...string) {
	c.borders = make([]string, 0, len(borders))
	for _, border := range borders {
		b := strings.TrimSpace(border)
		if len(b) > 0 {
			c.borders = append(c.borders, b)
		}
	}
}
