package rdbms

type demonymRecord struct {
	id         uint32
	countryId  uint16
	languageId uint16
	female     string
	male       string
}

func toDemonymRecord(scn scannable, result *demonymRecord) error {
	return scn.Scan(&result.id, &result.countryId, &result.languageId, &result.female, &result.male)
}

type demonymsDB struct {
	prepStmt prepStatementI
}

func (db *demonymsDB) createDemonyms(records []*demonymRecord) error {
	stmt, err := db.prepStmt.Prepare(`INSERT INTO demonyms (country_id, language_id, female, male) 
VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		return err
	}
	defer closeWithShowError(stmt)
	for _, record := range records {
		if err = stmt.QueryRow(record.countryId, record.languageId, record.female, record.male).Scan(&record.id); err != nil {
			return err
		}
	}
	return nil
}
