package rdbms

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type regionsDbTestSuite struct {
	databaseBaseTestSuite
	dtb regionDb
}

func (s *regionsDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = regionDb{prepStmt: s.db}
}
func (s *regionsDbTestSuite) Test_readOrCreateSubregion_transaction() {
	tests := []struct {
		continent string
		region    string
		subregion string
		want      RegionRecord
	}{
		{continent: "Europe", region: "Europe", subregion: "Western Europe",
			want: RegionRecord{RegionId: 3, ParentId: 2, RegionName: "Western Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Western Asia",
			want: RegionRecord{RegionId: 6, ParentId: 5, RegionName: "Western Asia"}},
		{continent: "Europe", region: "Europe", subregion: "Eastern Europe",
			want: RegionRecord{RegionId: 7, ParentId: 2, RegionName: "Eastern Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Eastern Asia",
			want: RegionRecord{RegionId: 8, ParentId: 5, RegionName: "Eastern Asia"}},
	}

	for i, tt := range tests {
		s.Run(fmt.Sprintf("%d_'%s'>'%s'>'%s'", i, tt.continent, tt.region, tt.subregion), func() {
			tx, err := s.databaseBaseTestSuite.db.Begin()
			s.Nil(err)
			s.NotNil(tx)
			defer func(tx *sql.Tx) {
				showError(tx.Rollback())
			}(tx)
			rdb := regionDb{prepStmt: tx}
			actual, err := rdb.readOrCreateSubregion(tt.continent, tt.region, tt.subregion)
			s.Nil(err)
			s.Equal(tt.want, actual)
			err = tx.Commit()
			s.Nil(err)
		})
	}

}
func (s *regionsDbTestSuite) Test_readOrCreateSubregion() {
	tests := []struct {
		continent string
		region    string
		subregion string
		want      RegionRecord
	}{
		{continent: "Europe", region: "Europe", subregion: "Western Europe",
			want: RegionRecord{RegionId: 3, ParentId: 2, RegionName: "Western Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Western Asia",
			want: RegionRecord{RegionId: 6, ParentId: 5, RegionName: "Western Asia"}},
		{continent: "Europe", region: "Europe", subregion: "Eastern Europe",
			want: RegionRecord{RegionId: 7, ParentId: 2, RegionName: "Eastern Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Eastern Asia",
			want: RegionRecord{RegionId: 8, ParentId: 5, RegionName: "Eastern Asia"}},
	}

	for i, tt := range tests {
		s.Run(fmt.Sprintf("%d_'%s'>'%s'>'%s'", i, tt.continent, tt.region, tt.subregion), func() {
			actual, err := s.dtb.readOrCreateSubregion(tt.continent, tt.region, tt.subregion)
			s.Nil(err)
			s.Equal(tt.want, actual)
		})
	}
}
func (s *regionsDbTestSuite) Test_readOrCreateSubregion_region_exists() {
	tests := []struct {
		continent string
		region    string
		subregion string
		want      RegionRecord
	}{
		{continent: "Europe", region: "Europe", subregion: "Western Europe",
			want: RegionRecord{RegionId: 3, ParentId: 2, RegionName: "Western Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Western Asia",
			want: RegionRecord{RegionId: 6, ParentId: 5, RegionName: "Western Asia"}},
		{continent: "Europe", region: "Europe", subregion: "Eastern Europe",
			want: RegionRecord{RegionId: 7, ParentId: 2, RegionName: "Eastern Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Eastern Asia",
			want: RegionRecord{RegionId: 8, ParentId: 5, RegionName: "Eastern Asia"}},
	}

	for _, tt := range tests {
		actual, err := s.dtb.readOrCreateSubregion(tt.continent, tt.region, tt.subregion)
		s.Nil(err)
		s.Equal(tt.want, actual)
	}
}
func (s *regionsDbTestSuite) Test_readOrCreateRegion_transaction() {
	tx, err := s.databaseBaseTestSuite.db.Begin()
	s.Nil(err)
	s.NotNil(tx)
	defer func(tx *sql.Tx) {
		showError(tx.Rollback())
	}(tx)
	rdb := regionDb{prepStmt: tx}
	name := "Europe"
	actual, err := rdb.readOrCreateRegion(name, name)
	s.Nil(err)
	err = tx.Commit()
	s.Nil(err)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateRegion() {
	name := "Europe"
	actual, err := s.dtb.readOrCreateRegion(name, name)
	s.Nil(err)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateRegion_continent_exists() {
	name := "Europe"
	_, err := s.dtb.readOrCreateContinents(name)
	s.Nil(err)
	actual, err := s.dtb.readOrCreateRegion(name, name)
	s.Nil(err)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateRegion_continent_duplicate() {
	name := "Europe"
	actual, err := s.dtb.readOrCreateRegion(name, name)
	s.Nil(err)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, actual)
	actual2, err2 := s.dtb.readOrCreateRegion(name, name)
	s.Nil(err2)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, actual2)
}
func (s *regionsDbTestSuite) Test_readOrCreateContinents_transaction() {
	tx, err := s.databaseBaseTestSuite.db.Begin()
	s.Nil(err)
	s.NotNil(tx)
	defer func(tx *sql.Tx) {
		showError(tx.Rollback())
	}(tx)
	rdb := regionDb{prepStmt: tx}
	names := []string{"Europe", "Asia"}
	actual, err := rdb.readOrCreateContinents(names...)
	s.Nil(err)
	err = tx.Commit()
	s.Nil(err)
	s.Equal([]RegionRecord{
		{RegionId: 1, ParentId: 0, RegionName: names[0]},
		{RegionId: 2, ParentId: 0, RegionName: names[1]}}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateContinents() {
	names := []string{"Europe", "Asia"}
	actual, err := s.dtb.readOrCreateContinents(names...)
	s.Nil(err)
	s.Equal([]RegionRecord{
		{RegionId: 1, ParentId: 0, RegionName: names[0]},
		{RegionId: 2, ParentId: 0, RegionName: names[1]}}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateContinents_one_exist() {
	names := []string{"Europe", "Asia"}
	actual, err := s.dtb.readOrCreateContinents(names[1])
	s.Nil(err)
	s.Equal([]RegionRecord{{RegionId: 1, ParentId: 0, RegionName: names[1]}}, actual)
	actual2, err2 := s.dtb.readOrCreateContinents(names...)
	s.Nil(err2)
	s.Equal([]RegionRecord{
		{RegionId: 2, ParentId: 0, RegionName: names[0]},
		{RegionId: 1, ParentId: 0, RegionName: names[1]}}, actual2)
}
func (s *regionsDbTestSuite) Test_readOrCreateContinents_duplicate() {
	names := []string{"Europe", "Asia"}
	actual, err := s.dtb.readOrCreateContinents(names...)
	s.Nil(err)
	expected := []RegionRecord{
		{RegionId: 1, ParentId: 0, RegionName: names[0]},
		{RegionId: 2, ParentId: 0, RegionName: names[1]}}
	s.Equal(expected, actual)
	actual2, err2 := s.dtb.readOrCreateContinents(names...)
	s.Nil(err2)
	s.Equal(expected, actual2)
}

func (s *regionsDbTestSuite) TestCreateContinent() {
	name := "Europe"
	reg, err := s.dtb.CreateContinent(name)
	s.Nil(err)
	s.Equal(RegionRecord{RegionId: 1, ParentId: 0, RegionName: name}, reg)
}
func (s *regionsDbTestSuite) TestCreateRegion() {
	name := "Europe"
	continent, err := s.dtb.CreateContinent(name)
	s.Nil(err)
	s.Equal(RegionRecord{RegionId: 1, ParentId: 0, RegionName: name}, continent)
	reg, errReg := s.dtb.CreateRegion(name, continent.RegionId)
	s.Nil(errReg)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, reg)
}
func (s *regionsDbTestSuite) TestCreateSubregion() {
	name := "Europe"
	continent, err := s.dtb.CreateContinent(name)
	s.Nil(err)
	reg, errReg := s.dtb.CreateRegion(name, continent.RegionId)
	s.Nil(errReg)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, reg)
	subregionName := "Western Europe"
	subregion, errSubreg := s.dtb.CreateRegion(subregionName, reg.RegionId)
	s.Nil(errSubreg)
	s.Equal(RegionRecord{RegionId: 3, ParentId: 2, RegionName: subregionName}, subregion)
}
func (s *regionsDbTestSuite) TestGetContinent() {
	name := "Europe"
	continent, errC := s.dtb.CreateContinent(name)
	s.Nil(errC)
	s.Equal(RegionRecord{RegionId: 1, ParentId: 0, RegionName: name}, continent)
	region, errR := s.dtb.CreateRegion(name, continent.RegionId)
	s.Nil(errR)
	subregionName := "Western Europe"
	_, errS := s.dtb.CreateRegion(subregionName, region.RegionId)
	s.Nil(errS)
	actual, err := s.dtb.GetContinent(name)
	s.Nil(err)
	s.Equal(continent, actual)
}
func (s *regionsDbTestSuite) TestGetRegion() {
	name := "Europe"
	continent, errC := s.dtb.CreateContinent(name)
	s.Nil(errC)
	region, errR := s.dtb.CreateRegion(name, continent.RegionId)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, region)
	s.Nil(errR)
	subregionName := "Western Europe"
	_, errS := s.dtb.CreateRegion(subregionName, region.RegionId)
	s.Nil(errS)
	actual, err := s.dtb.GetRegion(name)
	s.Nil(err)
	s.Equal(region, actual)
}
func (s *regionsDbTestSuite) TestGetSubregion() {
	name := "Europe"
	continent, errC := s.dtb.CreateContinent(name)
	s.Nil(errC)
	region, errR := s.dtb.CreateRegion(name, continent.RegionId)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, region)
	s.Nil(errR)
	subregionName := "Western Europe"
	subregion, errS := s.dtb.CreateRegion(subregionName, region.RegionId)
	s.Nil(errS)
	s.Equal(RegionRecord{RegionId: 3, ParentId: 2, RegionName: subregionName}, subregion)
	actual, err := s.dtb.GetSubregion(subregionName)
	s.Nil(err)
	s.Equal(subregion, actual)
}
func (s *regionsDbTestSuite) TestGetContinentNotFound() {
	name := "Europe"
	subregionName := "Western Europe"
	_, _, _, err := s.createRegions(name, name, subregionName)
	s.Nil(err)
	continent, errS := s.dtb.GetContinent(subregionName)
	s.NotNil(errS)
	s.True(errors.Is(errS, sql.ErrNoRows))
	s.Equal(RegionRecord{RegionId: 0, ParentId: 0, RegionName: subregionName}, continent)
}
func Test_regionsDbTestSuite(t *testing.T) {
	suite.Run(t, new(regionsDbTestSuite))
}
