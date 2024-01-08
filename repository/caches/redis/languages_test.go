package redis

import (
	"context"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type languageCacheTestSuit struct {
	redisTestSuit
	cache languageCache
}

func (s *languageCacheTestSuit) SetupSuite() {
	s.redisTestSuit.SetupSuite()
	s.cache = languageCache{rdb: s.rdb}
}

func Test_languagesDbTestSuite(t *testing.T) {
	suite.Run(t, new(languageCacheTestSuit))
}

func (s *languageCacheTestSuit) Test_readLanguageByCode_not_found() {
	actualId, actualName, err := s.cache.readLanguageByCode("none")
	s.ErrorIs(err, ErrNotFound)
	s.EqualValues(0, actualId)
	s.Equal("", actualName)
}
func (s *languageCacheTestSuit) Test_readLanguageByCode() {
	langId := uint16(1)
	langCode := "eng"
	langName := "English"
	s.NoError(s.cache.storeLanguage(langId, langCode, langName))

	actualId, actualName, err := s.cache.readLanguageByCode(langCode)
	s.Nil(err)
	s.Equal(langId, actualId)
	s.Equal(langName, actualName)
}
func (s *languageCacheTestSuit) Test_readLanguageById() {
	langId := uint16(1)
	langCode := "eng"
	langName := "English"
	s.NoError(s.cache.storeLanguage(langId, langCode, langName))
	actualCode, actualName, err := s.cache.readLanguageById(langId)
	s.Nil(err)
	s.Equal(langCode, actualCode)
	s.Equal(langName, actualName)
}
func (s *languageCacheTestSuit) Test_readLanguageById_not_found() {
	actualCode, actualName, err := s.cache.readLanguageById(123)
	s.ErrorIs(err, ErrNotFound)
	s.Equal("", actualCode)
	s.Equal("", actualName)
}
func (s *languageCacheTestSuit) Test_storeLanguage() {
	langId := uint16(1)
	langCode := "eng"
	langName := "English"
	s.NoError(s.cache.storeLanguage(langId, langCode, langName))
	n, err := s.cache.rdb.Exists(context.Background(), languageIdKey(langId)).Result()
	if err != nil {
		s.T().Fatal(err)
	}
	s.EqualValues(1, n)
	str, err := s.cache.rdb.HGet(context.Background(), languageIdKey(langId), code).Result()
	s.Nil(err)
	s.Equal(langCode, str)
	mp, err := s.rdb.HGetAll(context.Background(), languageIdKey(langId)).Result()
	s.Nil(err)
	s.Equal(langCode, mp[code])
	s.Equal(langName, mp[name])
	mp, err = s.rdb.HGetAll(context.Background(), languageCodeKey(langCode)).Result()
	s.Nil(err)
	s.Equal(strconv.Itoa(int(langId)), mp[id])
	str, err = s.rdb.HGet(context.Background(), languageCodeKey(langCode), id).Result()
	s.Nil(err)
	s.Equal(strconv.Itoa(int(langId)), str)
}
