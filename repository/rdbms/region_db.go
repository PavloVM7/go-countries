package rdbms

import "database/sql"

type regionDb struct {
	db *sql.DB
}

func (db *regionDb) GetSubregion(name string) (RegionRecord, error) {
	return db.getRegion(`SELECT s.region_id, s.parent_id, s.region_name
FROM regions r INNER JOIN (SELECT region_id FROM regions WHERE parent_id = 0) c ON r.parent_id = c.region_id 
    INNER JOIN regions s ON s.parent_id = r.region_id
WHERE s.region_name = $1`, name)
}
func (db *regionDb) GetRegion(name string) (RegionRecord, error) {
	return db.getRegion(`SELECT r.region_id, r.parent_id, r.region_name
FROM regions r INNER JOIN (SELECT region_id FROM regions WHERE parent_id = 0) c ON r.parent_id = c.region_id
WHERE r.region_name = $1`, name)
}
func (db *regionDb) GetContinent(name string) (RegionRecord, error) {
	return db.getRegion("SELECT * FROM regions WHERE parent_id = 0 AND region_name = $1", name)
}

func (db *regionDb) getRegion(sqlRequest, name string) (RegionRecord, error) {
	result := RegionRecord{RegionId: 0, ParentId: 0, RegionName: name}
	stmt, err := db.db.Prepare(sqlRequest)
	if err != nil {
		return result, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	err = stmt.QueryRow(result.RegionName).Scan(&result.RegionId, &result.ParentId, &result.RegionName)
	return result, err
}
func (db *regionDb) CreateContinent(name string) (RegionRecord, error) {
	return db.CreateRegion(name, 0)
}

func (db *regionDb) CreateRegion(name string, parentId uint32) (RegionRecord, error) {
	result := RegionRecord{RegionId: 0, ParentId: parentId, RegionName: name}
	stmt, err := db.db.Prepare("INSERT INTO regions (parent_id, region_name) VALUES ($1, $2) RETURNING region_id")
	if err != nil {
		return result, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	err = stmt.QueryRow(result.ParentId, result.RegionName).Scan(&result.RegionId)
	return result, err
}
