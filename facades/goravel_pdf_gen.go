package facades

import (
	"log"

	"goravel/packages/goravel_pdf_gen"
	"goravel/packages/goravel_pdf_gen/contracts"
)

func GoravelPdfGen() contracts.GoravelPdfGen {
	instance, err := goravel_pdf_gen.App.Make(goravel_pdf_gen.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.GoravelPdfGen)
}
