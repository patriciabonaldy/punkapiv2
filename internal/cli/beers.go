package cli

import (
	"fmt"

	service "github.com/patriciabonaldy/punkapiv2/internal/cli/fetching"
	"github.com/spf13/cobra"
)

// CobraFn function definion of run cobra command
type CobraFn func(cmd *cobra.Command, args []string)

const idFlag = "idFlag"
const nameFileFlag = "nameFile"

// InitBeersCmd initialize beers command
func InitBeersCmd(fetching service.Service) *cobra.Command {

	beersCmd := &cobra.Command{
		Use:   "beers",
		Short: "Print data about beers",
		Run:   runBeersFn(fetching),
	}

	beersCmd.Flags().StringP(idFlag, "i", "", "id de la beer")
	beersCmd.Flags().StringP(nameFileFlag, "n", "", "name of file")

	return beersCmd
}

func runBeersFn(fetching service.Service) CobraFn {
	return func(cmd *cobra.Command, args []string) {
		nameFile := ""
		id, _ := cmd.Flags().GetString(idFlag)
		var beers []Beer
		var err error

		beers, err = fetching.FetchBeers()

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(beers)

		}

	}
}
