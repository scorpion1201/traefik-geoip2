package traefik_geoip2

import (
	"context"
	"github.com/oschwald/geoip2-golang"
	"net"
	"net/http"
)

type Config struct {
	ASNFileName     string `json:"asnFileName,omitempty"`
	CityFileName    string `json:"cityFileName,omitempty"`
	CountryFileName string `json:"countryFileName,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		ASNFileName:     "/usr/share/GeoIP/GeoLite2-ASN.mmdb",
		CityFileName:    "/usr/share/GeoIP/GeoLite2-City.mmdb",
		CountryFileName: "/usr/share/GeoIP/GeoLite2-Country.mmdb",
	}
}

type GeoIP2 struct {
	next      http.Handler
	name      string
	AsnDB     *geoip2.Reader
	CityDB    *geoip2.Reader
	CountryDB *geoip2.Reader
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	asnDB, _ := geoip2.Open(config.ASNFileName)
	cityDB, _ := geoip2.Open(config.CityFileName)
	countryDB, _ := geoip2.Open(config.CountryFileName)

	return &GeoIP2{
		next:      next,
		name:      name,
		AsnDB:     asnDB,
		CityDB:    cityDB,
		CountryDB: countryDB,
	}, nil
}

func (g *GeoIP2) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		g.next.ServeHTTP(rw, req)
		return
	}

	// lookup provider
	if organization, err := g.getOrganization(net.ParseIP(host)); err == nil {
		rw.Header().Add("X-Edge-IPProvider", organization)
	}

	// lookup country
	if countryCode, err := g.getOrganization(net.ParseIP(host)); err == nil {
		rw.Header().Add("X-Edge-IPCountry", countryCode)
	}

	// lookup city
	if cityName, err := g.getOrganization(net.ParseIP(host)); err == nil {
		rw.Header().Add("X-Edge-IPCity", cityName)
	}

	g.next.ServeHTTP(rw, req)
}

func (g *GeoIP2) getOrganization(ipAddress net.IP) (string, error) {
	if g.AsnDB != nil {
		record, err := g.AsnDB.ASN(ipAddress)
		if err != nil {
			return "", err
		}
		return record.AutonomousSystemOrganization, nil
	}
	return "", nil
}

func (g *GeoIP2) getCountry(ipAddress net.IP) (string, error) {
	if g.CountryDB != nil {
		record, err := g.CountryDB.Country(ipAddress)
		if err != nil {
			return "", err
		}
		return record.Country.IsoCode, nil
	}
	return "", nil
}

func (g *GeoIP2) getCityName(ipAddress net.IP) (string, error) {
	if g.CityDB != nil {
		record, err := g.CityDB.City(ipAddress)
		if err != nil {
			return "", err
		}
		return record.City.Names["en"], nil
	}
	return "", nil
}
