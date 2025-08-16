package goravel_pdf_gen

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
	"goravel/packages/goravel_pdf_gen/routes"
)

const Binding = "goravel_pdf_gen"

var App foundation.Application

type ServiceProvider struct {
}

// Relationship returns the relationship of the service provider.
func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings:     []string{},
		Dependencies: []string{},
		ProvideFor:   []string{},
	}
}

// Register registers the service provider.
func (r *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return nil, nil
	})
}

// Boot boots the service provider, will be called after all service providers are registered.
func (r *ServiceProvider) Boot(app foundation.Application) {
	routes.Pdf(app)
	app.Publishes("./packages/goravel_pdf_gen", map[string]string{
		"setup/config/goravel_pdf_gen.go":                   app.ConfigPath("goravel_pdf_gen.go"),
		"databases/20250816105218_create_pdf_gens_table.go": app.DatabasePath("20250816105218_create_pdf_gens_table.go"),
		"assets/all.min.css":                                app.PublicPath("assets/all.min.css"),
		"assets/fa-solid-900.woff2":                         app.PublicPath("assets/fa-solid-900.woff2"),
		"assets/vue.global.js":                              app.PublicPath("assets/vue.global.js"),
	})
}
