package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	beer "github.com/patriciabonaldy/punkapiV2/internal"
	"github.com/patriciabonaldy/punkapiV2/internal/errors"
)

//AbvEndpoint endpoint abv
const (
	punkapiEndpoint = "https://api.punkapi.com/v2/beers"
	AbvEndpoint     = punkapiEndpoint + "?abv_gt="
)

//BeerRepo interface
type BeerRepo interface {
	GetBeers() ([]beer.Beer, error)
	GetBeersByComplemento(url string) ([]beer.Beer, error)
}

type repository struct {
	url string
}

// NewRepository initialize csv repository
func NewRepository() BeerRepo {
	return &repository{url: punkapiEndpoint}
}

// GetBeers fetch beers data from csv
func (r *repository) GetBeers() ([]beer.Beer, error) {
	var beers []beer.Beer
	response, err := http.Get(fmt.Sprintf("%v", r.url))

	if err != nil {
		return nil, errors.WrapUnreacheableBeerErr(err, "error obteniendo endpoint %v", r.url)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WrapUnreacheableBeerErr(err, "error leyendo el response %v", r.url)
	}

	err = json.Unmarshal(contents, &beers)
	if err != nil {
		return nil, errors.WrapUnreacheableBeerErr(err, "error parsing to beers")
	}

	return beers, nil
}

// GetBeers fetch beers data from csv
func (r *repository) GetBeersByComplemento(url string, complemento string) ([]beer.Beer, error) {
	var beers []beer.Beer
	response, err := http.Get(fmt.Sprintf("%v%v", url, complemento))

	if err != nil {
		return nil, errors.WrapUnreacheableBeerErr(err, "error obteniendo endpoint %v", r.url)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WrapUnreacheableBeerErr(err, "error leyendo el response %v", r.url)
	}

	err = json.Unmarshal(contents, &beers)
	if err != nil {
		return nil, errors.WrapUnreacheableBeerErr(err, "error parsing to beers")
	}

	return beers, nil
}
