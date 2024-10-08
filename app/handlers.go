package main

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RequestData struct {
	Name      string  `json:"name" binding:"required"`
	CreatedAt int64   `json:"createdAt" binding:"required"`
	UniqueId  *string `json:"uniqueId,omitempty"`
	Type      string  `json:"type,omitempty"`
}

func analyticsHandler(c *gin.Context) {
	lookback := 7
	if queryLookback := c.Query("lookback"); queryLookback != "" {
		if parsed, err := strconv.Atoi(queryLookback); err == nil && parsed > 0 {
			lookback = parsed
		}
	}

	redisKey := os.Getenv("REDIS_KEY")
	if redisKey == "" {
		redisKey = "analyticsEngine"
	}

	analyticsType := c.Query("type")
	if analyticsType != "" {
		redisKey += "-" + analyticsType
	}

	now := time.Now()
	dailyCutoff := now.AddDate(0, 0, -lookback).UnixMilli()
	weeklyCutoff := now.AddDate(0, 0, -(7 * lookback)).UnixMilli()
	monthlyCutoff := now.AddDate(0, -lookback, 0).UnixMilli()

	events, err := rdb.ZRangeByScore(ctx, redisKey, &redis.ZRangeBy{
		Min: strconv.FormatInt(monthlyCutoff, 10),
		Max: "+inf",
	}).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Failed to retrieve events from Redis: " + err.Error()})
		return
	}

	usages := make(map[string]map[string]map[string]int)
	global := map[string]map[string]int{
		"daily":   {},
		"weekly":  {},
		"monthly": {},
	}

	filterUniqueId := c.Query("uniqueId")

	for _, event := range events {
		var reqData RequestData
		if err := json.Unmarshal([]byte(event), &reqData); err != nil {
			continue
		}

		if filterUniqueId != "" && (reqData.UniqueId == nil || *reqData.UniqueId != filterUniqueId) {
			continue
		}

		if _, exists := usages[reqData.Name]; !exists {
			usages[reqData.Name] = map[string]map[string]int{
				"daily":   {},
				"weekly":  {},
				"monthly": {},
			}
		}

		createdAt := time.UnixMilli(reqData.CreatedAt)
		dateKey := createdAt.Format("2006-01-02")
		weekStart := createdAt.AddDate(0, 0, -int(createdAt.Weekday())).Format("2006-01-02")
		monthKey := createdAt.Format("2006-01")

		if reqData.CreatedAt >= dailyCutoff {
			usages[reqData.Name]["daily"][dateKey]++
			global["daily"][dateKey]++
		}

		if reqData.CreatedAt >= weeklyCutoff {
			usages[reqData.Name]["weekly"][weekStart]++
			global["weekly"][weekStart]++
		}

		if reqData.CreatedAt >= monthlyCutoff {
			usages[reqData.Name]["monthly"][monthKey]++
			global["monthly"][monthKey]++
		}
	}

	response := gin.H{
		"global": global,
		"usages": usages,
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": response})
}

func eventHandler(c *gin.Context) {
	var reqData RequestData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "Invalid input: " + err.Error()})
		return
	}

	redisKey := os.Getenv("REDIS_KEY")
	if redisKey == "" {
		redisKey = "analyticsEngine"
	}

	if reqData.Type != "" {
		redisKey += "-" + reqData.Type
	}

	value, err := json.Marshal(reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Failed to marshal request data."})
		return
	}

	err = rdb.ZAdd(ctx, redisKey, &redis.Z{
		Score:  float64(reqData.CreatedAt),
		Member: value,
	}).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Failed to save data in Redis: " + err.Error()})
		return
	}

	maxExpiryInDays := os.Getenv("MAX_AGE")
	maxExpiry, _ := strconv.Atoi(maxExpiryInDays)

	err = rdb.Expire(ctx, redisKey, time.Hour*24*time.Duration(maxExpiry)).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Failed to set expiration in Redis: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": "Event stored successfully!"})
}

func flushHandler(c *gin.Context) {
	redisKey := os.Getenv("REDIS_KEY")
	if redisKey == "" {
		redisKey = "analyticsEngine"
	}

	analyticsType := c.Query("type")
	if analyticsType != "" {
		redisKey += "-" + analyticsType
	}

	err := rdb.Del(ctx, redisKey).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": "Failed to flush data from Redis: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": "Data flushed successfully!"})
}

func statsHandler(c *gin.Context) {
	totalKeys, _ := rdb.DBSize(ctx).Result()
	memoryBytes := getMemoryUsage()

	uptime := getSystemUptime()

	stats := gin.H{
		"total_redis_keys":      totalKeys,
		"cpu_usage":             getCpuUsage(),
		"ram_usage":             formatBytes(memoryBytes),
		"ram_usage_bytes":       memoryBytes,
		"system_uptime_seconds": uptime,
		"system_uptime":         formatSystemUptime(uptime),
		"go_routines":           runtime.NumGoroutine(),
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": stats})
}

func infoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": "Analytics Engine is running."})
}

func notFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": "Route not found."})
}
