package rdbms

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
)

type Point struct {
	Lat float64
	Lng float64
}

func (p *Point) Scan(src interface{}) (err error) {
	var data []byte
	switch src := src.(type) {
	case []byte:
		data = src
	case string:
		data = []byte(src)
	case nil:
		return nil
	default:
		return errors.New("unsupported data type")
	}

	if len(data) == 0 {
		return nil
	}

	data = data[1 : len(data)-1] // drop the surrounding parentheses
	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			if p.Lat, err = strconv.ParseFloat(string(data[:i]), 64); err != nil {
				return err
			}
			if p.Lng, err = strconv.ParseFloat(string(data[i+1:]), 64); err != nil {
				return err
			}
			break
		}
	}
	return nil
}
func (p *Point) Value() (driver.Value, error) {
	out := []byte{'('}
	out = strconv.AppendFloat(out, p.Lat, 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, p.Lng, 'f', -1, 64)
	out = append(out, ')')
	return out, nil
}
func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.Lng, p.Lat)
}
