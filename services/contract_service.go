package contract

import (
	"bytes"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/hulutech-web/goravel_pdf_gen/models"
	"os"
	"path/filepath"
)

type ContractService struct {
}

func NewContractService() *ContractService {
	return &ContractService{}
}

// 创建合同并生成PDF
func (s *ContractService) CreateAndGeneratePDF(contract *models.Contract, templateContent string, variables models.ContractVariables) (string, error) {
	// 2. 生成合同编号
	contractNo := contract.GenerateContractNo()
	variables.ContractNo = contractNo
	contract.Content = contractNo // 可存储合同编号，或存储完整渲染后的HTML

	// 3. 渲染HTML模板
	renderedHTML, err := contract.RenderTemplate(templateContent, variables)
	if err != nil {
		return "", err
	}

	return renderedHTML, nil
}

// HTMLToPDF 将HTML内容转换为PDF字节流
func HTMLToPDF(htmlContent string) ([]byte, error) {
	// 1. 初始化PDF生成器
	pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("初始化PDF生成器失败: %w", err)
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

	pageOpts.NoImages.Set(false)     // 不禁用图片
	pageOpts.NoBackground.Set(false) // 不禁用背景（透明图片需要）

	// 3. 添加HTML内容作为PDF页面
	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(htmlContent)))

	pageOpts.UserStyleSheet.Set(`
		.sign-image { display: block !important; max-width: 180px !important; max-height: 80px !important; }
		.sign-image-container { height: 80px !important; }
	`) // 强制图片显示，避免CSS冲突导致隐藏
	page.PageOptions = pageOpts
	pdfGenerator.AddPage(page)
	// 4. 生成PDF
	if err = pdfGenerator.Create(); err != nil {
		// 打印调试日志，查看图片加载失败的具体原因（如Base64格式错误）
		return nil, fmt.Errorf("生成PDF失败: %w，调试日志: %s", err, pdfGenerator.LogLevel)
	}

	// 5. 返回PDF字节流
	return pdfGenerator.Bytes(), nil
}

// SavePDFToFile 将PDF字节流保存到文件（可选，用于测试或存储）
func SavePDFToFile(pdfBytes []byte, outputDir, filename string) (string, error) {
	// 创建输出目录（如不存在）
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 拼接完整路径
	filePath := filepath.Join(outputDir, filename)
	// 写入文件
	if err := os.WriteFile(filePath, pdfBytes, 0644); err != nil {
		return "", fmt.Errorf("保存PDF文件失败: %w", err)
	}

	return filePath, nil
}
