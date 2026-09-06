package utils

import (
	"Lancio/model"
	"fmt"
	"time"
)

// cacheKey 生成缓存 key
func cacheKey(query, lang string) string {
	return fmt.Sprintf("%s|%s", query, lang)
}

// SetResultCache 存入缓存（按 query+lang 作为 key）
func SetResultCache(query, lang string, r model.Result) {
	model.ResultCache.Mu.Lock()
	defer model.ResultCache.Mu.Unlock()
	model.ResultCache.Data[cacheKey(query, lang)] = model.CacheEntry{
		Result:  r,
		Expires: time.Now().Add(model.ResultCache.TTL),
	}
}

// GetResultCache 取缓存，未命中或过期返回 nil
func GetResultCache(query, lang string) *model.Result {
	model.ResultCache.Mu.RLock()
	entry, ok := model.ResultCache.Data[cacheKey(query, lang)]
	model.ResultCache.Mu.RUnlock()
	if !ok {
		return nil
	}
	if time.Now().After(entry.Expires) {
		// 已过期，异步删除避免阻塞读
		go func() {
			model.ResultCache.Mu.Lock()
			delete(model.ResultCache.Data, cacheKey(query, lang))
			model.ResultCache.Mu.Unlock()
		}()
		return nil
	}
	r := entry.Result
	return &r
}

// CleanupResultCache 清理所有过期条目
func CleanupResultCache() {
	model.ResultCache.Mu.Lock()
	defer model.ResultCache.Mu.Unlock()
	now := time.Now()
	for k, v := range model.ResultCache.Data {
		if now.After(v.Expires) {
			delete(model.ResultCache.Data, k)
		}
	}
}

// ClearResultCache 清空全部缓存
func ClearResultCache() {
	model.ResultCache.Mu.Lock()
	defer model.ResultCache.Mu.Unlock()
	model.ResultCache.Data = make(map[string]model.CacheEntry)
}
