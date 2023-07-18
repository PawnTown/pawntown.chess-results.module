package pawntown_chess_results_module

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type ChessResultsApiOptions struct {
	BaseUrl string
}

type Api struct {
	baseUrl string
}

func New(options *ChessResultsApiOptions) *Api {
	return &Api{
		baseUrl: "https://chess-results.com",
	}
}

func (api *Api) get(path string) (string, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", api.baseUrl, path))
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
