package controllers

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/goravel/framework/contracts/http"
	goravelfacades "github.com/goravel/framework/facades"
	"github.com/hulutech-web/goravel_pdf_gen/models"
)

type PDFController struct {
}

func NewPDFController() *PDFController {
	return &PDFController{}
}
func (c *PDFController) Generate(ctx http.Context) http.Response {
	var pdf models.PdfGen
	id := ctx.Request().RouteInt("id")
	goravelfacades.Orm().Query().Model(&models.PdfGen{}).Where("id=?", id).Find(&pdf)

	// 1. 初始化PDF生成器
	pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		ctx.Request().AbortWithStatusJson(http.StatusInternalServerError, http.Json{
			"msg": "初始化PDF生成器失败",
			"err": err.Error(),
		})
		return nil
	}

	// 2. 配置页面选项（支持中文、设置纸张大小等）
	pdfGenerator.Dpi.Set(600)
	pdfGenerator.NoCollate.Set(false)
	pdfGenerator.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfGenerator.MarginBottom.Set(40)
	pdfGenerator.MarginTop.Set(10)
	pdfGenerator.MarginLeft.Set(15)
	pdfGenerator.MarginRight.Set(15)

	//	// 3. 页面配置：强制启用所有资源加载
	pageOpts := wkhtmltopdf.NewPageOptions()
	// 允许所有类型的资源加载（关键修复）
	pageOpts.Allow.Unset()                     // 清除默认限制
	pageOpts.Allow.Set("*")                    // 允许所有URL和协议
	pageOpts.Allow.Set("*")                    // 允许所有URL
	pageOpts.DisableLocalFileAccess.Set(false) // 确保禁用本地文件访问的选项为false
	pageOpts.Encoding.Set("utf-8")

	pageOpts.NoImages.Set(false)     // 不禁用图片
	pageOpts.NoBackground.Set(false) // 不禁用背景（透明图片需要）

	// 3. 添加HTML内容作为PDF页面
	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(pdf.Html)))
	pageOpts.UserStyleSheet.Set(`
		.sign-image { display: block !important; max-width: 180px !important; max-height: 80px !important; }
		.sign-image-container { height: 80px !important; }
	`) // 强制图片显示，避免CSS冲突导致隐藏
	page.PageOptions = pageOpts
	pdfGenerator.AddPage(page)
	// 4. 生成PDF
	if err = pdfGenerator.Create(); err != nil {
		ctx.Request().AbortWithStatusJson(http.StatusInternalServerError, http.Json{
			"msg": "初始化PDF生成器失败",
			"err": err.Error(),
		})
		return nil
	}

	// 6. 返回PDF文件流
	return ctx.Response().
		Header("Content-Type", "application/pdf; charset=utf-8").
		Header("Content-Disposition", "attachment; filename*=UTF-8''生成的文档.pdf").
		Data(http.StatusOK, "application/pdf", pdfGenerator.Bytes())
}

// 获取随机默认模板
func (c *PDFController) GetDefaultTemplate(ctx http.Context) http.Response {
	randomTpl := models.PdfGen{}
	goravelfacades.Orm().Query().InRandomOrder().First(&randomTpl)
	return ctx.Response().Success().Json(http.Json{"data": randomTpl})
}

// 保存模板
func (c *PDFController) SaveTemplate(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")
	params := ctx.Request().Input("params")
	goravelfacades.Orm().Query().Model(&models.PdfGen{}).Where("id=?", id).Update("params", params)
	return ctx.Response().Success().Json(http.Json{"msg": "保存成功"})
}

// 保存模板
func (c *PDFController) SaveHTML(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")
	params := ctx.Request().Input("html")
	goravelfacades.Orm().Query().Model(&models.PdfGen{}).Where("id=?", id).Update("html", params)
	return ctx.Response().Success().Json(http.Json{"msg": "保存成功"})
}

// 获取所有模板列表
func (c *PDFController) GetIndexTemplate(ctx http.Context) http.Response {
	tpls := []models.PdfGen{}
	goravelfacades.Orm().Query().InRandomOrder().First(&tpls)
	return ctx.Response().Success().Json(http.Json{"data": tpls})
}
