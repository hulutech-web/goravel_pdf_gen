package databases

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250816105218CreatePdfGensTable struct{}

// Signature The unique signature for the migration.
func (r *M20250816105218CreatePdfGensTable) Signature() string {
	return "20250816105218_create_pdf_gens_table"
}

// Up Run the migrations.
func (r *M20250816105218CreatePdfGensTable) Up() error {
	if !facades.Schema().HasTable("pdf_gens") {
		return facades.Schema().Create("pdf_gens", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.Json("params")
			table.Text("html")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250816105218CreatePdfGensTable) Down() error {
	return facades.Schema().DropIfExists("pdf_gens")
}
