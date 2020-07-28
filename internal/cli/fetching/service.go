package fetching

import (
	"fmt"
	"math"

	"github.com/patriciabonaldy/punkapiV2/internal/errors"
	beerscli "github.com/patriciabonaldy/punkapiv2/internal"
	storage "github.com/patriciabonaldy/punkapiv2/internal/cli/storage"
)

// Service provides beer fetching operations
type Service interface {
	// FetchBeers fetch all beers from repository
	FetchBeers() ([]beerscli.Beer, error)
	// FetchByID filter all beers and get only the beer that match with given id
	FetchByID(id int) (beerscli.Beer, error)
}

//service
type service struct {
	bR storage.BeerRepo
}

// NewFetchy initialize csv repository
func NewFetchy(r storage.BeerRepo) Service {
	return &service{r}
}

func findBeersbyAbv(r *service, chanelAbv *chan float32, chanelBeer *chan []Beer) {
	var abv float32
	for {
		abv <- chanelAbv
		beers, err := r.bR.GetBeersByComplemento(AbvEndpoint, abv)
		if err == nil {
			beerscli.ChanelBeer <- beers
		}
	}
}

func (r *service) FetchBeers() ([]beerscli.Beer, error) {
	//chanelAbv y chanelBeer  channel
	var chanelAbv = make(chan float32, 10)
	var chanelBeer = make(chan []Beer, 10)
	beers, err := r.bR.GetBeers()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, beer := range beers {
		chanelAbv <- beer.Abv
		go findBeersbyAbv(&r, &chanelAbv, &chanelBeer)
	}

	return beers, nil
}

func (r *service) FetchByID(id int) (beerscli.Beer, error) {
	beers, err := r.FetchBeers()

	if err != nil {
		return beerscli.Beer{}, err
	}

	beersPerRoutine := 10
	numRoutines := numOfRoutines(len(beers), beersPerRoutine)

	b := make(chan beerscli.Beer)
	done := make(chan bool, numRoutines)

	for i := 0; i < numRoutines; i++ {
		toSearch := make([]beerscli.Beer, beersPerRoutine)
		copy(toSearch[:], beers[i:i+beersPerRoutine])

		go func(beers []beerscli.Beer, b chan beerscli.Beer, done chan bool) {
			for _, beer := range beers {
				if beer.ProductID == id {
					b <- beer
				}
			}
			done <- true
		}(toSearch, b, done)
	}

	var beer beerscli.Beer
	i := 0
	for i < numRoutines {
		select {
		case beer = <-b:
			return beer, nil
		case <-done:
			i++
		}
	}

	return beerscli.Beer{}, errors.NewUnreacheableBeerErr("No existe la beer con id %v", id)
}

func numOfRoutines(numOfBeers, beersPerRoutine int) int {
	return int(math.Ceil(float64(numOfBeers) / float64(beersPerRoutine)))
}
