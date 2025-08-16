package config

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"
)

func init() {
	config := facades.Config()
	config.Add("goravel_pdf_gen", map[string]any{
		"contracts": map[string]any{
			"name":            "string",
			"phone":           "string",
			"address":         "string",
			"contract_number": "string",
			"contract_date":   "string",
			"contract_amount": "string",
			"contract_type":   "string",
			"contract_status": "string",
			"contract_remark": "string",
		},
	})

	/*方法路径*/
	config.Add("pdf_save_path", path.Public("pdf"))
	config.Add("pdf_prefix", "pdf_prefix")
	config.Add("getDefaultTemplate", "getDefaultTemplate")
	config.Add("saveTemplate", "saveTemplate")
	config.Add("saveHTML", "saveHTML")
	config.Add("getIndexTemplate", "getIndexTemplate")
	config.Add("generate", "generate")
}
