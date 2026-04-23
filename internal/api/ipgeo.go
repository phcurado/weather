package api

import (
	"errors"
	"net/url"
)

const DefaultIPGeoBase = "https://ipwho.is"

// ErrIPGeoFailed is returned when IP geolocation has no usable result.
var ErrIPGeoFailed = errors.New("ip geolocation failed")

type ipGeoResp struct {
	Success   bool    `json:"success"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  struct {
		ID string `json:"id"`
	} `json:"timezone"`
}

// LocateByIP resolves the caller's approximate coordinates from their public IP.
func (c *Client) LocateByIP() (Coords, error) {
	var out ipGeoResp
	if err := c.getJSON(c.ipGeoBase+"/", url.Values{}, &out); err != nil {
		return Coords{}, err
	}
	if !out.Success || (out.Latitude == 0 && out.Longitude == 0) {
		return Coords{}, ErrIPGeoFailed
	}
	return Coords{
		Name:     out.City,
		Country:  out.Country,
		Lat:      out.Latitude,
		Lon:      out.Longitude,
		Timezone: out.Timezone.ID,
	}, nil
}
