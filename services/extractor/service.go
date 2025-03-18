package extractor

import (
	"context"
	"log"
)

var ()

func Start(ctx context.Context) {

	<-ctx.Done()
	log.Println("Stopping extractor..")
}
