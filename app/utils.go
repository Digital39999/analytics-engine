package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/cpu"
)

var startTime = time.Now()

func loadEnvVars() (string, string, error) {
	_ = godotenv.Load()

	redisURL := os.Getenv("REDIS_URL")
	port := os.Getenv("PORT")

	if redisURL == "" || port == "" {
		return "", "", errors.New("missing REDIS_URL or PORT in environment variables")
	}

	if os.Getenv("API_AUTH") == "" {
		return "", "", errors.New("missing API_AUTH in environment variables")
	}

	if os.Getenv("MAX_AGE") == "" {
		return "", "", errors.New("missing MAX_AGE in environment variables")
	}

	return redisURL, port, nil
}

func getCpuUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return 0
	}

	if len(percent) > 0 {
		return math.Round(percent[0]*100) / 100
	}

	return 0
}

func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	totalAllocated := m.Alloc + m.TotalAlloc
	return totalAllocated
}

func formatBytes(bytes uint64) string {
	const (
		_         = iota
		KB uint64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2fTB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

func getSystemUptime() int64 {
	uptime := time.Since(startTime).Seconds()
	return int64(uptime)
}

func formatSystemUptime(uptime int64) string {
	months := uptime / 2592000
	uptime %= 2592000
	days := uptime / 86400
	uptime %= 86400
	hours := uptime / 3600
	uptime %= 3600
	minutes := uptime / 60
	seconds := uptime % 60

	result := ""
	if months > 0 {
		result += fmt.Sprintf("%dm ", months)
	}
	if days > 0 {
		result += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%ds", seconds)
	}

	return result
}
