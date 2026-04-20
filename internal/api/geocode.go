package api

import "net/url"

type geocodeResp struct {
	Results []Coords `json:"results"`
}

// Geocode resolves a city name to coordinates using Open-Meteo's geocoder.
// Returns ErrCityNotFound when no results match.
func (c *Client) Geocode(name string) (Coords, error) {
	q := url.Values{}
	q.Set("name", name)
	q.Set("count", "1")
	q.Set("language", "en")
	q.Set("format", "json")

	var out geocodeResp
	if err := c.getJSON(c.geocodeBase+"/v1/search", q, &out); err != nil {
		return Coords{}, err
	}
	if len(out.Results) == 0 {
		return Coords{}, ErrCityNotFound
	}
	return out.Results[0], nil
}
