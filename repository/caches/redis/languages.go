package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

const (
	id   = "id"
	code = "code"
	name = "name"
)

func languageIdKey(languageId uint16) string {
	return "language:" + strconv.Itoa(int(languageId))
}
func languageCodeKey(languageCode string) string {
	return "language:code:" + languageCode
}

type languageCache struct {
	rdb redis.Cmdable
}

func (cache *languageCache) readLanguageByCode(languageCode string) (langId uint16, langName string, err error) {
	ctx := context.Background()
	codeKey := languageCodeKey(languageCode)
	var strId string
	strId, err = cache.rdb.HGet(ctx, codeKey, id).Result()
	wrapErr := func(er error) {
		err = wrapLanguageErr(er, codeKey, languageCode, langName)
	}
	if errors.Is(err, redis.Nil) {
		wrapErr(ErrNotFound)
		return
	} else if err != nil {
		wrapErr(err)
	}
	if strId == "" {
		wrapErr(ErrNotFound)
		return
	}
	var lId int64
	lId, err = asInt(strId)
	if err != nil {
		wrapErr(err)
		return
	}
	langId = uint16(lId)
	idKey := languageIdKey(langId)
	langName, err = cache.rdb.HGet(ctx, idKey, name).Result()
	if err != nil {
		wrapErr(err)
		return
	}
	return
}

func (cache *languageCache) readLanguageById(languageId uint16) (langCode, langName string, err error) {
	ctx := context.Background()
	idKey := languageIdKey(languageId)
	wrapErr := func(er error) {
		err = wrapLanguageErr(er, idKey, langCode, langName)
	}
	cmd := cache.rdb.HMGet(ctx, idKey, code, name)
	if cmd.Err() != nil {
		wrapErr(cmd.Err())
		return
	}

	var res struct {
		Code string `redis:"code"`
		Name string `redis:"name"`
	}
	err = cmd.Scan(&res)
	if err != nil {
		wrapErr(err)
		return
	}
	langCode = res.Code
	langName = res.Name
	if langCode == "" && langName == "" {
		wrapErr(ErrNotFound)
	}
	return
}

func (cache *languageCache) storeLanguage(languageId uint16, languageCode, languageName string) error {
	pipe := cache.rdb.Pipeline()
	cxt := context.Background()
	idKey := languageIdKey(languageId)
	wrapErr := func(err error) error { return wrapLanguageErr(err, idKey, languageCode, languageName) }
	_, err := pipe.HSet(cxt, languageIdKey(languageId), code, languageCode, name, languageName).Result()
	if err != nil {
		return wrapErr(err)
	}
	_, err = pipe.HSet(cxt, languageCodeKey(languageCode), id, strconv.Itoa(int(languageId))).Result()
	if err != nil {
		return wrapErr(err)
	}
	_, err = pipe.Exec(cxt)
	if err != nil {
		return wrapErr(err)
	}
	return nil
}

func wrapLanguageErr(err error, key, code, name string) error {
	return fmt.Errorf("%w; %s:%s:'%s'", err, key, code, name)
}
