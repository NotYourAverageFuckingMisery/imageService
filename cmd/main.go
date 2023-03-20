package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/NotYourAverageFuckingMisery/imageService/internal/service"
	"github.com/NotYourAverageFuckingMisery/imageService/internal/store"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()

	store := store.NewImageStore("./imageStore")
	server := service.NewServer(store, 10, 100)

	go func() {
		runtime.SetCPUProfileRate(500)
		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

	}()

	server.Run("0.0.0.0:5051", "0.0.0.0:5069")

}
