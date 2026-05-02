package geo

import (
	"fmt"

	"github.com/gookit/slog"
	"github.com/ip2location/ip2location-go/v9"
)

type Service struct {
	db *ip2location.DB
}

func New(path string) (*Service, error) {
	db, err := ip2location.OpenDB(path)
	if err != nil {
		return nil, err
	}

	return &Service{db: db}, nil
}

func (s *Service) Close() {
	s.db.Close()
}

func isLocalIP(ip string) bool {
	return ip == "127.0.0.1" || ip == "::1"
}

func (s *Service) GetGeoString(ip string) (country string, city string) {
	if isLocalIP(ip) {
		return "Local", "Local"
	}

	res, err := s.db.Get_all(ip)
	if err != nil {
		slog.Errorf("Failed to get geolocation data: %v", err)
	}
	if s.db == nil {
		panic("DB is nil")
	}
	fmt.Println("DB pointer:", s.db)

	fmt.Println("Country:", res.Country_long)
	fmt.Println("City:", res.City)
	fmt.Println("Timezone:", res.Timezone)

	slog.Info(res)

	return res.Country_long, res.City
}
