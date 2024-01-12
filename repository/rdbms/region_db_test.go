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
		want      regionRecord
	}{
		{continent: "Europe", region: "Europe", subregion: "Western Europe",
			want: regionRecord{regionId: 3, parentId: 2, regionName: "Western Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Western Asia",
			want: regionRecord{regionId: 6, parentId: 5, regionName: "Western Asia"}},
		{continent: "Europe", region: "Europe", subregion: "Eastern Europe",
			want: regionRecord{regionId: 7, parentId: 2, regionName: "Eastern Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Eastern Asia",
			want: regionRecord{regionId: 8, parentId: 5, regionName: "Eastern Asia"}},
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
		want      regionRecord
	}{
		{continent: "Europe", region: "Europe", subregion: "Western Europe",
			want: regionRecord{regionId: 3, parentId: 2, regionName: "Western Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Western Asia",
			want: regionRecord{regionId: 6, parentId: 5, regionName: "Western Asia"}},
		{continent: "Europe", region: "Europe", subregion: "Eastern Europe",
			want: regionRecord{regionId: 7, parentId: 2, regionName: "Eastern Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Eastern Asia",
			want: regionRecord{regionId: 8, parentId: 5, regionName: "Eastern Asia"}},
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
		want      regionRecord
	}{
		{continent: "Europe", region: "Europe", subregion: "Western Europe",
			want: regionRecord{regionId: 3, parentId: 2, regionName: "Western Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Western Asia",
			want: regionRecord{regionId: 6, parentId: 5, regionName: "Western Asia"}},
		{continent: "Europe", region: "Europe", subregion: "Eastern Europe",
			want: regionRecord{regionId: 7, parentId: 2, regionName: "Eastern Europe"}},
		{continent: "Asia", region: "Asia", subregion: "Eastern Asia",
			want: regionRecord{regionId: 8, parentId: 5, regionName: "Eastern Asia"}},
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
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateRegion() {
	name := "Europe"
	actual, err := s.dtb.readOrCreateRegion(name, name)
	s.Nil(err)
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateRegion_continent_exists() {
	name := "Europe"
	_, err := s.dtb.readOrCreateContinents(name)
	s.Nil(err)
	actual, err := s.dtb.readOrCreateRegion(name, name)
	s.Nil(err)
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateRegion_continent_duplicate() {
	name := "Europe"
	actual, err := s.dtb.readOrCreateRegion(name, name)
	s.Nil(err)
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, actual)
	actual2, err2 := s.dtb.readOrCreateRegion(name, name)
	s.Nil(err2)
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, actual2)
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
	s.Equal([]regionRecord{
		{regionId: 1, parentId: 0, regionName: names[0]},
		{regionId: 2, parentId: 0, regionName: names[1]}}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateContinents() {
	names := []string{"Europe", "Asia"}
	actual, err := s.dtb.readOrCreateContinents(names...)
	s.Nil(err)
	s.Equal([]regionRecord{
		{regionId: 1, parentId: 0, regionName: names[0]},
		{regionId: 2, parentId: 0, regionName: names[1]}}, actual)
}
func (s *regionsDbTestSuite) Test_readOrCreateContinents_one_exist() {
	names := []string{"Europe", "Asia"}
	actual, err := s.dtb.readOrCreateContinents(names[1])
	s.Nil(err)
	s.Equal([]regionRecord{{regionId: 1, parentId: 0, regionName: names[1]}}, actual)
	actual2, err2 := s.dtb.readOrCreateContinents(names...)
	s.Nil(err2)
	s.Equal([]regionRecord{
		{regionId: 2, parentId: 0, regionName: names[0]},
		{regionId: 1, parentId: 0, regionName: names[1]}}, actual2)
}
func (s *regionsDbTestSuite) Test_readOrCreateContinents_duplicate() {
	names := []string{"Europe", "Asia"}
	actual, err := s.dtb.readOrCreateContinents(names...)
	s.Nil(err)
	expected := []regionRecord{
		{regionId: 1, parentId: 0, regionName: names[0]},
		{regionId: 2, parentId: 0, regionName: names[1]}}
	s.Equal(expected, actual)
	actual2, err2 := s.dtb.readOrCreateContinents(names...)
	s.Nil(err2)
	s.Equal(expected, actual2)
}
func (s *regionsDbTestSuite) Test_readRegionByIds() {
	names := []string{"Europe", "Asia"}
	actual, err := s.dtb.readOrCreateContinents(names...)
	s.Nil(err)
	expected := []regionRecord{
		{regionId: 1, parentId: 0, regionName: names[0]},
		{regionId: 2, parentId: 0, regionName: names[1]}}
	s.Equal(expected, actual)
	continents, errR := s.dtb.readRegionsByIds(1, 2)
	s.Nil(errR)
	s.Equal(expected, continents)
}
func (s *regionsDbTestSuite) TestCreateContinent() {
	name := "Europe"
	reg, err := s.dtb.CreateContinent(name)
	s.Nil(err)
	s.Equal(regionRecord{regionId: 1, parentId: 0, regionName: name}, reg)
}
func (s *regionsDbTestSuite) TestCreateRegion() {
	name := "Europe"
	continent, err := s.dtb.CreateContinent(name)
	s.Nil(err)
	s.Equal(regionRecord{regionId: 1, parentId: 0, regionName: name}, continent)
	reg, errReg := s.dtb.CreateRegion(name, continent.regionId)
	s.Nil(errReg)
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, reg)
}
func (s *regionsDbTestSuite) TestCreateSubregion() {
	name := "Europe"
	continent, err := s.dtb.CreateContinent(name)
	s.Nil(err)
	reg, errReg := s.dtb.CreateRegion(name, continent.regionId)
	s.Nil(errReg)
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, reg)
	subregionName := "Western Europe"
	subregion, errSubreg := s.dtb.CreateRegion(subregionName, reg.regionId)
	s.Nil(errSubreg)
	s.Equal(regionRecord{regionId: 3, parentId: 2, regionName: subregionName}, subregion)
}
func (s *regionsDbTestSuite) TestGetContinent() {
	name := "Europe"
	continent, errC := s.dtb.CreateContinent(name)
	s.Nil(errC)
	s.Equal(regionRecord{regionId: 1, parentId: 0, regionName: name}, continent)
	region, errR := s.dtb.CreateRegion(name, continent.regionId)
	s.Nil(errR)
	subregionName := "Western Europe"
	_, errS := s.dtb.CreateRegion(subregionName, region.regionId)
	s.Nil(errS)
	actual, err := s.dtb.getContinent(name)
	s.Nil(err)
	s.Equal(continent, actual)
}
func (s *regionsDbTestSuite) TestGetRegion() {
	name := "Europe"
	continent, errC := s.dtb.CreateContinent(name)
	s.Nil(errC)
	region, errR := s.dtb.CreateRegion(name, continent.regionId)
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, region)
	s.Nil(errR)
	subregionName := "Western Europe"
	_, errS := s.dtb.CreateRegion(subregionName, region.regionId)
	s.Nil(errS)
	actual, err := s.dtb.GetRegion(name)
	s.Nil(err)
	s.Equal(region, actual)
}
func (s *regionsDbTestSuite) TestGetSubregion() {
	name := "Europe"
	continent, errC := s.dtb.CreateContinent(name)
	s.Nil(errC)
	region, errR := s.dtb.CreateRegion(name, continent.regionId)
	s.Equal(regionRecord{regionId: 2, parentId: 1, regionName: name}, region)
	s.Nil(errR)
	subregionName := "Western Europe"
	subregion, errS := s.dtb.CreateRegion(subregionName, region.regionId)
	s.Nil(errS)
	s.Equal(regionRecord{regionId: 3, parentId: 2, regionName: subregionName}, subregion)
	actual, err := s.dtb.getSubregion(subregionName)
	s.Nil(err)
	s.Equal(subregion, actual)
}
func (s *regionsDbTestSuite) TestGetContinentNotFound() {
	name := "Europe"
	subregionName := "Western Europe"
	_, _, _, err := s.createRegions(name, name, subregionName)
	s.Nil(err)
	continent, errS := s.dtb.getContinent(subregionName)
	s.NotNil(errS)
	s.True(errors.Is(errS, sql.ErrNoRows))
	s.Equal(regionRecord{regionId: 0, parentId: 0, regionName: subregionName}, continent)
}
func Test_regionsDbTestSuite(t *testing.T) {
	suite.Run(t, new(regionsDbTestSuite))
}
