package cache

import "time"

const (
	CacheTtl5s  = time.Second * 5
	CacheTtl10s = time.Second * 10
	CacheTtl15s = time.Second * 15
	CacheTtl30s = time.Second * 30

	CacheTtl1m  = time.Minute
	CacheTtl2m  = time.Minute * 2
	CacheTtl5m  = time.Minute * 5
	CacheTtl10m = time.Minute * 10
	CacheTtl15m = time.Minute * 15
	CacheTtl30m = time.Minute * 30

	CacheTtl1h  = time.Hour
	CacheTtl3h  = time.Hour * 3
	CacheTtl6h  = time.Hour * 6
	CacheTtl9h  = time.Hour * 9
	CacheTtl12h = time.Hour * 12

	CacheTtl1d  = time.Hour * 24
	CacheTtl2d  = CacheTtl1d * 2
	CacheTtl3d  = CacheTtl1d * 3
	CacheTtl4d  = CacheTtl1d * 4
	CacheTtl7d  = CacheTtl1d * 7
	CacheTtl9d  = CacheTtl1d * 9
	CacheTtl10d = CacheTtl1d * 10
	CacheTtl15d = CacheTtl1d * 15
	CacheTtl30d = CacheTtl1d * 30
)
