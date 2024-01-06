package rdbms

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
)

type RegionRecord struct {
	RegionId   uint32
	ParentId   uint32
	RegionName string
}

func toRegionRecord(scn scannable, result *RegionRecord) error {
	return scn.Scan(&result.RegionId, &result.ParentId, &result.RegionName)
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
)

type regionDb struct {
	prepStmt prepStatementI
}

func (db *regionDb) readOrCreateSubregion(continent, region, subregion string) (RegionRecord, error) {
	var (
		sub        RegionRecord
		stmtSelect *sql.Stmt
		err        error
	)
	sub.RegionName = subregion
	stmtSelect, err = db.prepStmt.Prepare(selectSubregionByNames)
	if err != nil {
		return sub, err
	}
	defer closeAndShowError(stmtSelect)
	row := stmtSelect.QueryRow(region, subregion)
	err = toRegionRecord(row, &sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			reg, errR := db.readOrCreateRegion(continent, region)
			if errR != nil {
				err = fmt.Errorf("couldn't get/create region '%s/%s', '%w'", region, continent, errR)
			}
			sub.ParentId = reg.RegionId
			err = db.insertRegion(&sub)
		}
	}
	return sub, err
}
func (db *regionDb) readOrCreateRegion(continent, region string) (RegionRecord, error) {
	var (
		reg        RegionRecord
		stmtSelect *sql.Stmt
		err        error
	)
	reg.RegionName = region
	stmtSelect, err = db.prepStmt.Prepare(selectRegionByNames)
	if err != nil {
		return reg, err
	}
	defer closeAndShowError(stmtSelect)
	row := stmtSelect.QueryRow(continent, region)
	err = toRegionRecord(row, &reg)
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
			reg.ParentId = continents[0].RegionId
			err = db.insertRegion(&reg)
		}
	}
	return reg, err
}
func (db *regionDb) readOrCreateContinents(continents ...string) ([]RegionRecord, error) {
	var (
		stmtSelect *sql.Stmt
		stmtInsert *sql.Stmt
		err        error
	)

	stmtSelect, err = db.prepStmt.Prepare(selectContinent)
	if err != nil {
		return nil, err
	}
	defer closeAndShowError(stmtSelect)
	stmtInsert, err = db.prepStmt.Prepare(insertRegion)
	if err != nil {
		return nil, err
	}
	defer closeAndShowError(stmtInsert)
	result := lists.NewLinkedList[RegionRecord]()
	for _, continent := range continents {
		row := stmtSelect.QueryRow(continent)
		reg := RegionRecord{}
		err = toRegionRecord(row, &reg)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				reg.RegionName = continent
				reg.ParentId = 0
				err = db.insertQueryAndScanRegion(stmtInsert, &reg)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		result.AddLast(reg)
	}
	return result.ToArray(), nil
}
func (db *regionDb) GetSubregion(name string) (RegionRecord, error) {
	return db.getRegion(selectSubregion, name)
}
func (db *regionDb) GetRegion(name string) (RegionRecord, error) {
	return db.getRegion(selectRegion, name)
}
func (db *regionDb) GetContinent(name string) (RegionRecord, error) {
	return db.getRegion(selectContinent, name)
}

func (db *regionDb) getRegion(sqlRequest, name string) (RegionRecord, error) {
	result := RegionRecord{RegionId: 0, ParentId: 0, RegionName: name}
	stmt, err := db.prepStmt.Prepare(sqlRequest)
	if err != nil {
		return result, err
	}
	defer closeAndShowError(stmt)
	row := stmt.QueryRow(result.RegionName)
	err = toRegionRecord(row, &result)
	return result, err
}
func (db *regionDb) CreateContinent(name string) (RegionRecord, error) {
	return db.CreateRegion(name, 0)
}

func (db *regionDb) CreateRegion(name string, parentId uint32) (RegionRecord, error) {
	result := RegionRecord{RegionId: 0, ParentId: parentId, RegionName: name}
	stmt, err := db.prepStmt.Prepare(insertRegion)
	if err != nil {
		return result, err
	}
	defer closeAndShowError(stmt)
	err = db.insertQueryAndScanRegion(stmt, &result)
	return result, err
}
func (db *regionDb) insertRegion(record *RegionRecord) error {
	stmtInsert, err := db.prepStmt.Prepare(insertRegion)
	if err != nil {
		return err
	}
	defer closeAndShowError(stmtInsert)
	return db.insertQueryAndScanRegion(stmtInsert, record)
}
func (db *regionDb) insertQueryAndScanRegion(insertStmt *sql.Stmt, record *RegionRecord) error {
	return insertStmt.QueryRow(record.ParentId, record.RegionName).Scan(&record.RegionId)
}
