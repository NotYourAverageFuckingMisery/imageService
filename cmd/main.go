package main

import (
	"github.com/NotYourAverageFuckingMisery/imageService/internal/service"
	"github.com/NotYourAverageFuckingMisery/imageService/internal/store"
	// v1 "github.com/NotYourAverageFuckingMisery/imageService/proto/v1"
)

func main() {

	store := store.NewImageStore("./imageStore")
	server := service.NewServer(store, 10, 100)

	server.Run("0.0.0.0:5051", "0.0.0.0:5069")

}
