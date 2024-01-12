package rdbms

import (
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"github.com/lib/pq"
	"math"
	"pm.com/go-countries/domain"
)

type capitalInfoRecord struct {
	capitalId uint32
	point     domain.LatLng
}

func toCapitalInfoRecord(row scannable, record *capitalInfoRecord) error {
	var point Point
	err := row.Scan(&record.capitalId, &point)
	if err != nil {
		return err
	}
	record.point = domain.LatLng{Lat: float32(point.Lat), Lng: float32(point.Lng)}
	return nil
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
