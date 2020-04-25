package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const apiURL = "https://api.e-hentai.org/api.php"
const exURL = "https://exhentai.org"

type ExHentaiAPI struct {
	client *http.Client
}

func NewExHentaiAPI(cookies []*http.Cookie) (*ExHentaiAPI, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(exURL)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}

	return &ExHentaiAPI{client: client}, nil
}

func (ex *ExHentaiAPI) RequestGalleryMetadata(rawurl string) (*GalleryMetadata, error) {
	galleryURL, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	b, err := newAPIRequest(*galleryURL)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", apiURL, b)
	req.Header.Set("Content-Type", "application/json")

	res, err := ex.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	apiRes := &apiResponse{}
	if err := json.NewDecoder(res.Body).Decode(apiRes); err != nil {
		return nil, err
	}
	if apiRes.GMetaData[0].Error != "" {
		return nil, errors.New(apiRes.GMetaData[0].Error)
	}

	return &apiRes.GMetaData[0], nil
}

func (ex *ExHentaiAPI) SearchGalleries(searchQuery string, page int) ([]*GalleryMetadata, error) {
	searchURL := BuildQuery(
		[]string{"f_search", "page"},
		[]string{searchQuery, strconv.Itoa(page)},
	)

	req, _ := http.NewRequest("GET", searchURL.String(), nil)
	res, err := ex.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var gms []*GalleryMetadata

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find(".gl1t").Each(func(i int, s *goquery.Selection) {
		rawurl, _ := s.Find("a").Attr("href")
		gm, err := ex.RequestGalleryMetadata(rawurl)
		if err == nil {
			gms = append(gms, gm)
		}
	})

	return gms, nil
}

func (ex *ExHentaiAPI) NumSearchGalleries(searchQuery string) (n int, err error) {
	searchURL := BuildQuery(
		[]string{"f_search"},
		[]string{searchQuery},
	)

	req, _ := http.NewRequest("GET", searchURL.String(), nil)
	res, err := ex.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	s := doc.Find(".ip")
	if s.Text() != "" {
		tokens := strings.Split(s.Text(), " ")
		n, err = strconv.Atoi(strings.ReplaceAll(tokens[1], ",", ""))
	}

	return
}
