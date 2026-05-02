package geo

import (
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

func (s *Service) GetGeoString(ip string) (countryCode string, country string, city string) {
	if isLocalIP(ip) {
		return "LC", "Local", "Local"
	}

	res, err := s.db.Get_all(ip)
	if err != nil {
		slog.Errorf("Failed to get geolocation data: %v", err)
		return "Unknown", "Unknown", "Unknown"
	}

	if res.Country_long == "" || res.Country_long == "This parameter is unavailable..." {
		return "Unknown", "Unknown", "Unknown"
	}

	return res.Country_short, res.Country_long, res.City
}
