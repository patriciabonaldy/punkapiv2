package main

import (
	"os"
	"runtime/pprof"

	cli "github.com/patriciabonaldy/punkapiv2/internal/cli/storage"

	service "github.com/patriciabonaldy/punkapiv2/internal/cli/fetching"
	storage "github.com/patriciabonaldy/punkapiv2/internal/cli/storage"
	"github.com/spf13/cobra"
)

func main() {
	//CPU profiling code start here
	f, _ := os.Create("beers.cpu.prof")
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	repo := storage.NewRepository()
	fetch := service.NewFetchy(repo)
	rootCmd := &cobra.Command{Use: "beers-cli"}
	rootCmd.AddCommand(cli.InitBeersCmd(fetch))
	rootCmd.Execute()

	//Memory profiling code start here
	f2, _ := os.Create("beers.mem.prof")
	defer f2.Close()
	pprof.WriteHeapProfile(f2)

}
