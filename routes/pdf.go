package routes

import (
	"bytes"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"goravel/packages/goravel_pdf_gen/controllers"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func Pdf(app foundation.Application) {
	router := app.MakeRoute()
	//渲染web端页面
	router.Get("/pdf_design", func(ctx http.Context) http.Response {
		//导入当前路径上一层中的dashboard.html文件
		// 获取当前源文件的绝对路径
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			// 处理错误
			panic("无法获取当前文件路径")
		}
		// 当前文件所在目录
		currentDir := filepath.Dir(filename)
		// 向上一层目录
		parentDir := filepath.Dir(currentDir)
		// 目标文件路径
		filePath := filepath.Join(parentDir, "dashboard.html")
		// 打开文件
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		var buf bytes.Buffer

		// 将文件内容复制到内存缓冲区
		_, err = io.Copy(&buf, file)
		if err != nil {
			panic(err)
		}
		// 现在可以通过 buf.Bytes() 或 buf.String() 访问内容
		content := buf.String()
		return ctx.Response().View().Make("design.tmpl", map[string]any{
			"content": template.HTML(content),
		})
	})

	// web端页面数据交互
	router.Prefix(facades.Config().GetString("pdf_prefix")).Group(func(r route.Router) {
		/*此处添加PDF业务路由*/
		pdfCtrl := controllers.NewPDFController()
		r.Get("getDefaultTemplate", pdfCtrl.GetDefaultTemplate)
		r.Post("saveTemplate/{id}", pdfCtrl.SaveTemplate)
		r.Post("saveHTML/{id}", pdfCtrl.SaveHTML)
		r.Get("getIndexTemplate", pdfCtrl.GetIndexTemplate)
		r.Get("generate/{id}", pdfCtrl.Generate)
	})
}
