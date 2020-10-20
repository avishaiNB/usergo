// Package doublecache cache provides advanced high performance cache services.
package doublecache

import (
	"sync"
	"time"

	cache "github.com/thelotter-enterprise/usergo/core/cache"
)

const (
	defaultBackupExpirationFactor float64 = 2.0
)

type lockKey struct {
	Region, Key string
}

// DoubleCache is a cache with 2-layer of cache
type DoubleCache struct {
	primary                *cache.Cache
	backup                 *cache.Cache
	locks                  map[lockKey]chan bool
	mu                     sync.RWMutex
	BackupExpirationFactor float64
}

// Config is used to configure the double cache
type Config struct {
	Expiration             time.Duration
	Timeout                time.Duration
	JitterFactor           float64
	BackupExpirationFactor float64
}

// DefaultConfig returns a default configuration of the double cache
func DefaultConfig() Config {
	return Config{
		Expiration:             cache.DefaultExpiration,
		Timeout:                cache.DefaultTimeout,
		JitterFactor:           cache.DefaultJitterFactor,
		BackupExpirationFactor: defaultBackupExpirationFactor,
	}
}

// NewDoubleCache create a DoubleCache
func NewDoubleCache(config Config) *DoubleCache {
	return &DoubleCache{
		primary: cache.NewCache(cache.Config{
			Expiration:   config.Expiration,
			Timeout:      config.Timeout,
			JitterFactor: config.JitterFactor,
		}),
		backup: cache.NewCache(cache.Config{
			Expiration:   time.Duration(float64(config.Expiration) * config.BackupExpirationFactor),
			Timeout:      config.Timeout,
			JitterFactor: config.JitterFactor,
		}),
		locks:                  make(map[lockKey]chan bool),
		BackupExpirationFactor: config.BackupExpirationFactor,
	}
}

// Get an item from the primary or backup cache
func (d *DoubleCache) Get(region, key string) (interface{}, bool) {
	v, err := d.primary.Get(region, key)
	if err == nil {
		return v, true
	}
	v, err = d.backup.Get(region, key)
	if err == nil {
		return v, true
	}
	return nil, false
}

// GetDefault will get an item from the default region
func (d *DoubleCache) GetDefault(key string) (interface{}, bool) {
	return d.Get(cache.DefaultRegion, key)
}

// GetOrCreate get an item
// If the information expired in the primary, but exist in the backup,
// the refresh function will run on the background and will set the return value in the
func (d *DoubleCache) GetOrCreate(
	region, key string,
	expiration time.Duration,
	refreshFunction cache.RefreshFunction,
) (interface{}, error) {
	v, err := d.primary.Get(region, key)
	if err == nil {
		return v, nil
	}

	backupValue, err := d.backup.Get(region, key)
	if err == nil {
		lock := d.getLock(region, key)
		select {
		case lock <- true:
			go func() {
				_, _ = d.refreshData(region, key, expiration, refreshFunction)
				<-lock
			}()
			return backupValue, nil
		default:
			return backupValue, nil
		}
	}
	v, err = d.refreshData(region, key, expiration, refreshFunction)
	return v, err
}

// GetOrCreateDefault will return an item
// from primary cache, backup cache or refresh function following this order in the default region
// If the information expired in the primary, but exist in the backup, the refresh function will run on the background
func (d *DoubleCache) GetOrCreateDefault(
	key string,
	expiration time.Duration,
	refreshFunction cache.RefreshFunction,
) (interface{}, error) {
	return d.GetOrCreate(cache.DefaultRegion, key, expiration, refreshFunction)
}

func (d *DoubleCache) refreshData(
	region, key string,
	expiration time.Duration,
	refreshFunction cache.RefreshFunction,
) (interface{}, error) {
	v, err := d.primary.GetOrCreate(region, key, expiration, refreshFunction)
	if err == nil {
		d.backup.Set(region, key, v, expiration)
		return v, nil
	}
	return nil, err
}

func (d *DoubleCache) getLock(region, key string) chan bool {
	d.mu.Lock()
	lock, found := d.locks[lockKey{region, key}]
	if !found {
		lock = make(chan bool, 1)
		d.locks[lockKey{region, key}] = lock
	}
	d.mu.Unlock()
	return lock
}

// Set will assign the value to the { region, key } in primary and backup
func (d *DoubleCache) Set(region, key string, value interface{}, expiration time.Duration) {
	d.primary.Set(region, key, value, expiration)
	d.backup.Set(region, key, value, expiration*2)
}

// SetDefault will assign the value to the key in the default region in primary and backup
func (d *DoubleCache) SetDefault(key string, value interface{}, expiration time.Duration) {
	d.primary.SetDefault(key, value, expiration)
	d.backup.SetDefault(key, value, expiration*2)
}

// InvalidateRegion will remove all entries in the region of the primary and backup cache
func (d *DoubleCache) InvalidateRegion(region string) {
	_ = d.primary.InvalidateRegion(region)
	_ = d.backup.InvalidateRegion(region)
}
