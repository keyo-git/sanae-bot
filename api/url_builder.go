package api

import (
	"fmt"
	"math"
	"net/url"
)

func BuildQuery(params, values []string) *url.URL {
	u, _ := url.Parse(exURL)
	n := int(math.Min(float64(len(params)), float64(len(values))))

	parameters := url.Values{}
	for i := 0; i < n; i++ {
		parameters.Add(params[i], values[i])
	}

	u.RawQuery = parameters.Encode()
	return u
}

func BuildURLFromMetadata(gm *GalleryMetadata) string {
	return exURL + fmt.Sprintf("/g/%d/%s/", gm.Gid, gm.Token)
}
