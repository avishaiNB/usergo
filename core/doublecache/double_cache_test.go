package doublecache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	cache "github.com/thelotter-enterprise/usergo/core/cache"
)

func Test_DoubleCache(t *testing.T) {
	tests := map[string]struct {
		region     string
		key        string
		expiration time.Duration
		input      interface{}
		action     func(t *testing.T, doubleCache *DoubleCache, regionName string, key string, input interface{}, expiration time.Duration)
		pValue     interface{}
		pError     error
		bValue     interface{}
		bError     error
	}{
		"Set": {
			region:     "1",
			key:        "1",
			expiration: 1 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)
			},
			pValue: 3,
			pError: nil,
			bValue: 3,
			bError: nil,
		},
		"Set default": {
			region:     cache.DefaultRegion,
			key:        "1",
			expiration: 1 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.SetDefault(key, input, expiration)
			},
			pValue: 3,
			pError: nil,
			bValue: 3,
			bError: nil,
		},
		"Get from primary": {
			region:     "1",
			key:        "1",
			expiration: 2 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				fn := func() (interface{}, error) {
					return 4, nil
				}
				v, err := doubleCache.GetOrCreate(region, key, cache.DefaultExpiration, fn)

				assert.Equal(t, input, v)
				assert.Equal(t, nil, err)
			},
			pValue: 3,
			pError: nil,
			bValue: 3,
			bError: nil,
		},
		"Value expire in primary, Get from backup": {
			region:     "2",
			key:        "2",
			expiration: 20 * time.Millisecond,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				// We wait for expiration * (1 + jitter)
				wait := time.Duration((1 + doubleCache.primary.JitterFactor) * float64(expiration))
				time.Sleep(wait + time.Millisecond)
			},
			pValue: nil,
			pError: cache.ErrKeyNotFound,
			bValue: 3,
			bError: nil,
		},
		"Value expire in primary, Get from backup and call backend to renew data": {
			region:     "2",
			key:        "2",
			expiration: 1 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				// We wait for expiration * (1 + jitter)
				wait := time.Duration((1 + doubleCache.primary.JitterFactor) * float64(expiration))
				time.Sleep(wait + time.Millisecond)

				fn := func() (interface{}, error) {
					return 4, nil
				}
				v, err := doubleCache.GetOrCreate(region, key, expiration, fn)
				assert.Equal(t, 3, v)
				assert.Equal(t, nil, err)
				time.Sleep(500 * time.Millisecond)
			},
			pValue: 4,
			pError: nil,
			bValue: 4,
			bError: nil,
		},
		"Value expire in primary, Get from backup and call backend to renew data with timeout": {
			region:     "2",
			key:        "2",
			expiration: 1 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				// We wait for expiration * (1 + jitter)
				wait := time.Duration((1 + doubleCache.primary.JitterFactor) * float64(expiration))
				time.Sleep(wait + time.Millisecond)

				fn := func() (interface{}, error) {
					time.Sleep(doubleCache.backup.Timeout + 1)
					return 4, nil
				}
				_, _ = doubleCache.GetOrCreate(region, key, expiration, fn)
			},
			pValue: nil,
			pError: cache.ErrKeyNotFound,
			bValue: 3,
			bError: nil,
		},
		"Value missing in both": {
			region:     "2",
			key:        "2",
			expiration: 2 * time.Second,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				fn := func() (interface{}, error) {
					return 4, nil
				}
				v, err := doubleCache.GetOrCreate(region, key, expiration, fn)
				assert.Equal(t, 4, v)
				assert.Equal(t, nil, err)
			},
			pValue: 4,
			pError: nil,
			bValue: 4,
			bError: nil,
		},
		"Value expire in both": {
			region:     "2",
			key:        "2",
			expiration: 20 * time.Millisecond,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, key, input, expiration)

				// We wait until the item is expire in both backups
				// Expiration time * (backup factor + jitter factor )
				wait := time.Duration(float64(expiration) * (doubleCache.BackupExpirationFactor + doubleCache.backup.JitterFactor))
				time.Sleep(wait + time.Second)
			},
			pValue: nil,
			pError: cache.ErrKeyNotFound,
			bValue: nil,
			bError: cache.ErrKeyNotFound,
		},
		"Region not found": {
			region:     "2",
			key:        "2",
			expiration: 2 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
			},
			pValue: nil,
			pError: cache.ErrRegionKeyNotFound,
			bValue: nil,
			bError: cache.ErrRegionKeyNotFound,
		},
		"Item not found": {
			region:     "2",
			key:        "2",
			expiration: 2 * time.Second,
			input:      3,
			action: func(t *testing.T, doubleCache *DoubleCache, region, key string, input interface{}, expiration time.Duration) {
				doubleCache.Set(region, "fake", input, expiration)
			},
			pValue: nil,
			pError: cache.ErrKeyNotFound,
			bValue: nil,
			bError: cache.ErrKeyNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			config := DefaultConfig()
			config.Timeout = 500 * time.Millisecond

			c := NewDoubleCache(config)
			tc.action(t, c, tc.region, tc.key, tc.input, tc.expiration)

			v, err := c.primary.Get(tc.region, tc.key)
			assert.Equal(t, tc.pValue, v)
			assert.Equal(t, tc.pError, err)

			v, err = c.backup.Get(tc.region, tc.key)
			assert.Equal(t, tc.bValue, v)
			assert.Equal(t, tc.bError, err)
		})
	}
}
