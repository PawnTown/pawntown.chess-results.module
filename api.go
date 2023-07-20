package pawntown_chess_results_module

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ChessResultsApiOptions struct {
	BaseUrl string
}

type Api struct {
	client  *http.Client
	baseUrl string
}

func New(options *ChessResultsApiOptions) *Api {

	return &Api{
		client:  &http.Client{},
		baseUrl: "https://chess-results.com",
	}
}

func (api *Api) get(path string) (string, error) {
	return api.req("GET", path, nil)
}

func (api *Api) post(path string, values url.Values) (string, error) {
	return api.req("POST", path, values)
}

func (api *Api) req(method string, path string, values url.Values) (string, error) {
	url := fmt.Sprintf("%s/%s", api.baseUrl, path)

	var body io.Reader = nil
	if values != nil {
		body = strings.NewReader(values.Encode())
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	if values != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := api.client.Do(req)
	if err != nil {
		return "", err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(resBody), nil
}

func (api *Api) Tournament(id string) (*Tournament, error) {
	return newTournament(api, id)
}
