package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"time"
)

func CacheMiddleware(cache Cache, expiration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		cacheKey := generateCacheKey(c.Request.URL.Path, c.Request.URL.RawQuery)

		cachedResponse, err := cache.Get(ctx, cacheKey)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Failed to retrieve cache"})
			return
		}

		if cachedResponse != nil {
			c.Data(200, "application/json", []byte(cachedResponse.(string)))
			c.Abort()
			return
		}

		writer := &responseWriter{body: make([]byte, 0), ResponseWriter: c.Writer}
		c.Writer = writer
		c.Next()

		if c.Writer.Status() == 200 {
			_ = cache.Set(ctx, cacheKey, string(writer.body), expiration)
		}
	}
}

func generateCacheKey(path, query string) string {
	hash := sha256.New()
	hash.Write([]byte(path + "?" + query))
	return hex.EncodeToString(hash.Sum(nil))
}

type responseWriter struct {
	gin.ResponseWriter
	body []byte
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body = append(rw.body, data...)
	return rw.ResponseWriter.Write(data)
}
