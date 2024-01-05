package rdbms

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type regionsDbTestSuite struct {
	databaseBaseTestSuite
	dtb regionDb
}

func (s *regionsDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = regionDb{db: s.db}
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

func Test_regionsDbTestSuite(t *testing.T) {
	suite.Run(t, new(regionsDbTestSuite))
}
