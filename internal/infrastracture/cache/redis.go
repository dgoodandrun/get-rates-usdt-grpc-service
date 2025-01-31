package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type GeoServiceProxy struct {
	geoService service.GeoService
	cache      *redis.Client
}

func NewGeoServiceProxy(geoService service.GeoService, cache *redis.Client) *GeoServiceProxy {
	return &GeoServiceProxy{
		geoService: geoService,
		cache:      cache,
	}
}

func (g *GeoServiceProxy) AddressSearch(input string) ([]*models.Address, error) {
	//start := time.Now()
	//defer func() {
	//	metrics.CacheDuration.
	//		WithLabelValues("AddressSearch").
	//		Observe(time.Since(start).Seconds())
	//}()

	ctx := context.Background()
	cacheKey := fmt.Sprintf("address_search:%s", input)

	cachedData, err := g.cache.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		addresses, err := g.geoService.AddressSearch(input)
		if err != nil {
			return nil, err
		}

		data, err := Serialize(addresses)
		if err != nil {
			return nil, err
		}
		g.cache.Set(ctx, cacheKey, data, 10*time.Minute)

		return addresses, nil
	} else if err != nil {
		return nil, err
	}

	var addresses []*models.Address
	if err = Deserialize([]byte(cachedData), &addresses); err != nil {
		return nil, err
	}
	log.Printf("Get cached data: %v", addresses)

	return addresses, nil
}

func (g *GeoServiceProxy) GeoCode(lat, lng string) ([]*models.Address, error) {
	//start := time.Now()
	//defer func() {
	//	metrics.CacheDuration.
	//		WithLabelValues("GeoCode").
	//		Observe(time.Since(start).Seconds())
	//}()

	ctx := context.Background()
	cacheKey := fmt.Sprintf("geocode:%s,%s", lat, lng)

	cachedData, err := g.cache.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		addresses, err := g.geoService.GeoCode(lat, lng)
		if err != nil {
			return nil, err
		}

		data, err := Serialize(addresses)
		if err != nil {
			return nil, err
		}
		g.cache.Set(ctx, cacheKey, data, 10*time.Minute)

		return addresses, nil
	} else if err != nil {
		return nil, err
	}

	var addresses []*models.Address
	if err = Deserialize([]byte(cachedData), &addresses); err != nil {
		return nil, err
	}
	log.Printf("Get cached data: %v", addresses)

	return addresses, nil
}
