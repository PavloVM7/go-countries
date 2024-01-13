package rdbms

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"github.com/lib/pq"
)

type regionRecord struct {
	regionId   uint32
	parentId   uint32
	regionName string
}

func toRegionRecord(scn scannable, result *regionRecord) error {
	return scn.Scan(&result.regionId, &result.parentId, &result.regionName)
}

const (
	selectContinent = "SELECT * FROM regions WHERE parent_id = 0 AND region_name = $1"

	selectRegion = `SELECT r.region_id, r.parent_id, r.region_name
FROM regions r INNER JOIN (SELECT region_id FROM regions WHERE parent_id = 0) c ON r.parent_id = c.region_id
WHERE r.region_name = $1`

	selectRegionByNames = `SELECT r.region_id, r.parent_id, r.region_name
FROM regions r, regions c 
WHERE r.region_name = $2 AND r.parent_id = c.region_id AND c.parent_id = 0 AND c.region_name = $1`

	selectSubregion = `SELECT s.region_id, s.parent_id, s.region_name FROM regions s, regions r, regions c
WHERE s.region_name = $1 AND s.parent_id = r.region_id AND r.parent_id = c.region_id AND c.parent_id = 0`

	selectSubregionByNames = `SELECT s.region_id, s.parent_id, s.region_name FROM regions s, regions r, regions c  
WHERE s.region_name = $2 AND s.parent_id = r.region_id AND r.region_name = $1 AND r.parent_id = c.region_id AND c.region_id = 0`

	insertRegion = "INSERT INTO regions (parent_id, region_name) VALUES ($1, $2) RETURNING region_id"

	selectRegionByIds = `SELECT region_id, parent_id, region_name FROM regions WHERE region_id = ANY ($1)`
)

type regionDb struct {
	prepStmt prepStatementI
}

func (db *regionDb) readOrCreateSubregion(continent, region, subregion string) (regionRecord, error) {
	var (
		sub        regionRecord
		stmtSelect *sql.Stmt
		err        error
	)
	sub.regionName = subregion
	stmtSelect, err = db.prepStmt.Prepare(selectSubregionByNames)
	if err != nil {
		return sub, err
	}
	defer closeWithShowError(stmtSelect)
	selSubregion := func(sub *regionRecord) error {
		return toRegionRecord(stmtSelect.QueryRow(region, subregion), sub)
	}
	err = selSubregion(&sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			reg, errR := db.readOrCreateRegion(continent, region)
			if errR != nil {
				err = fmt.Errorf("couldn't get/create region '%s/%s', '%w'", region, continent, errR)
			}
			sub.parentId = reg.regionId
			err = db.insertRegion(&sub)
			if err != nil && isErrorUniqueViolation(err) {
				err = selSubregion(&sub)
			}
		}
	}
	return sub, err
}
func (db *regionDb) readOrCreateRegion(continent, region string) (regionRecord, error) {
	var (
		reg        regionRecord
		stmtSelect *sql.Stmt
		err        error
	)
	reg.regionName = region
	stmtSelect, err = db.prepStmt.Prepare(selectRegionByNames)
	if err != nil {
		return reg, err
	}
	defer closeWithShowError(stmtSelect)
	selRegion := func(reg *regionRecord) error {
		return toRegionRecord(stmtSelect.QueryRow(continent, region), reg)
	}
	err = selRegion(&reg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			continents, errC := db.readOrCreateContinents(continent)
			if errC != nil {
				err = fmt.Errorf("couldn't get/create continent '%s', '%w'", continent, errC)
			} else if len(continents) == 0 {
				err = fmt.Errorf("couldn't get/create continent '%s'", continent)
			}
			if err != nil {
				return reg, err
			}
			reg.parentId = continents[0].regionId
			err = db.insertRegion(&reg)
			if err != nil && isErrorUniqueViolation(err) {
				err = selRegion(&reg)
			}
		}
	}
	return reg, err
}
func (db *regionDb) readOrCreateContinents(continents ...string) ([]regionRecord, error) {
	var (
		stmtSelect *sql.Stmt
		stmtInsert *sql.Stmt
		err        error
	)

	stmtSelect, err = db.prepStmt.Prepare(selectContinent)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmtSelect)
	stmtInsert, err = db.prepStmt.Prepare(insertRegion)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmtInsert)
	selContinent := func(reg *regionRecord) error {
		return toRegionRecord(stmtSelect.QueryRow(reg.regionName), reg)
	}
	result := lists.NewLinkedList[regionRecord]()
	for _, continent := range continents {
		reg := regionRecord{parentId: 0, regionName: continent}
		err = selContinent(&reg)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			err = insertRegionRecord(stmtInsert, &reg)
			if err != nil && isErrorUniqueViolation(err) {
				err = selContinent(&reg)
			}
			if err != nil {
				return nil, err
			}
		}
		result.AddLast(reg)
	}
	return result.ToArray(), nil
}
func (db *regionDb) getSubregion(name string) (regionRecord, error) {
	return db.getRegionRecord(selectSubregion, name)
}
func (db *regionDb) GetRegion(name string) (regionRecord, error) {
	return db.getRegionRecord(selectRegion, name)
}
func (db *regionDb) getContinent(name string) (regionRecord, error) {
	return db.getRegionRecord(selectContinent, name)
}
func (db *regionDb) readRegionsByIds(ids ...uint32) ([]regionRecord, error) {
	stmtSelect, err := db.prepStmt.Prepare(selectRegionByIds)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmtSelect)
	rows, er := stmtSelect.Query(pq.Array(ids))
	if er != nil {
		return nil, er
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[regionRecord]()
	for rows.Next() {
		reg := regionRecord{}
		if err = toRegionRecord(rows, &reg); err != nil {
			return nil, err
		}
		result.AddLast(reg)
	}
	return result.ToArray(), nil
}
func (db *regionDb) getRegionRecord(sqlRequest, name string) (regionRecord, error) {
	result := regionRecord{regionId: 0, parentId: 0, regionName: name}
	stmt, err := db.prepStmt.Prepare(sqlRequest)
	if err != nil {
		return result, err
	}
	defer closeWithShowError(stmt)
	row := stmt.QueryRow(result.regionName)
	err = toRegionRecord(row, &result)
	return result, err
}
func (db *regionDb) CreateContinent(name string) (regionRecord, error) {
	return db.CreateRegion(name, 0)
}

func (db *regionDb) CreateRegion(name string, parentId uint32) (regionRecord, error) {
	result := regionRecord{regionId: 0, parentId: parentId, regionName: name}
	stmt, err := db.prepStmt.Prepare(insertRegion)
	if err != nil {
		return result, err
	}
	defer closeWithShowError(stmt)
	err = insertRegionRecord(stmt, &result)
	return result, err
}
func (db *regionDb) insertRegion(record *regionRecord) error {
	stmtInsert, err := db.prepStmt.Prepare(insertRegion)
	if err != nil {
		return err
	}
	defer closeWithShowError(stmtInsert)
	return insertRegionRecord(stmtInsert, record)
}
func insertRegionRecord(insertStmt *sql.Stmt, record *regionRecord) error {
	return insertStmt.QueryRow(record.parentId, record.regionName).Scan(&record.regionId)
}
