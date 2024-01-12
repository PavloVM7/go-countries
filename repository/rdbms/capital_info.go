package rdbms

import (
	"errors"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"github.com/lib/pq"
	"math"
	"pm.com/go-countries/domain"
	"strconv"
)

type capitalInfoRecord struct {
	capitalId uint32
	point     domain.LatLng
}

func toCapitalInfoRecord(row scannable, record *capitalInfoRecord) error {
	var data []byte
	err := row.Scan(&record.capitalId, &data)
	if err != nil {
		return err
	}
	lat, lng, err := parseToLatLng(data)
	if err != nil {
		return err
	}
	record.point = domain.LatLng{Lat: float32(lat), Lng: float32(lng)}
	return nil
}

// ToDo: create struct Pointer{float64, float64} for this
func parseToLatLng(src []byte) (lat float64, lng float64, err error) {
	if len(src) == 0 {
		err = errors.New("empty data")
		return
	}
	data := src[1 : len(src)-1] // drops the surrounding parentheses
	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			if lat, err = strconv.ParseFloat(string(data[:i]), 64); err != nil {
				return
			}
			if lng, err = strconv.ParseFloat(string(data[i+1:]), 64); err != nil {
				return
			}
			break
		}
	}
	return
}

type capitalInfoDb struct {
	prepStmt prepStatementI
}

func (db *capitalInfoDb) createCapitalsInfo(capitalIds []uint32, points []domain.LatLng) error {
	stmt, err := db.prepStmt.Prepare("INSERT INTO country_capital_info VALUES ($1, POINT($2, $3))")
	if err != nil {
		return err
	}
	round := func(f float32) float64 {
		n := 100_000.
		return math.Round(float64(f)*n) / n
	}
	defer closeWithShowError(stmt)
	minLen := min(len(capitalIds), len(points))
	for i := 0; i < minLen; i++ {
		_, err = stmt.Exec(capitalIds[i], round(points[i].Lat), round(points[i].Lng))
		if err != nil {
			return err
		}
	}
	return nil
}
func (db *capitalInfoDb) readCapitalInfo(capitalIds ...uint32) ([]capitalInfoRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT capital_id, location FROM country_capital_info WHERE capital_id = ANY($1)")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	rows, er := stmt.Query(pq.Array(capitalIds))
	if er != nil {
		return nil, er
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[capitalInfoRecord]()
	for rows.Next() {
		reg := capitalInfoRecord{}
		if err = toCapitalInfoRecord(rows, &reg); err != nil {
			return nil, err
		}
		result.AddLast(reg)
	}
	return result.ToArray(), nil
}
