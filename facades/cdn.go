package facades

import (
	"log"

	"github.com/goravel-kit/cdn"
	"github.com/goravel-kit/cdn/contracts"
)

func Cdn() contracts.Cdn {
	instance, err := cdn.App.Make(cdn.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Cdn)
}
