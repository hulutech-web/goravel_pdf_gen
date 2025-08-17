<p align="center">
  <img src="https://github.com/hulutech-web/goravel_pdf_gen/blob/master/images/logo.png?raw=true" width="900"  />
</p>

<p align="center">
  <img src="https://github.com/hulutech-web/goravel_pdf_gen/blob/master/images/img2.png?raw=true" width="900"  />
</p>

<p align="center">
  <img src="https://github.com/hulutech-web/goravel_pdf_gen/blob/master/images/img3.png?raw=true" width="900"  />
</p>


# PDF 生成扩展包 for Goravel

## 项目简介

这是一个基于 Goravel 框架的 PDF 生成扩展包，提供强大的 HTML 模板编辑和 PDF 生成功能。扩展包集成了 Monaco 编辑器提供专业的代码编辑体验，并使用 wkhtmltopdf 实现高质量的 HTML 转 PDF 功能。

## 技术特点

- **前后端分离架构**：前端使用 Vue.js + Monaco Editor，后端基于 Goravel
- **专业代码编辑**：集成 Monaco Editor 提供代码高亮、智能提示和格式化功能
- **实时预览**：编辑 HTML 模板时可实时查看渲染效果
- **动态变量替换**：通过模板引擎动态修改合同内容
- **JSON Schema 支持**：使用 JSON Schema 定义表单结构和验证规则
- **配置化路由**：后端路由可通过配置文件灵活配置

## 功能特性

1. **可视化表单设计**：基于 JSON Schema 动态生成表单
2. **专业模板编辑**：
    - Monaco Editor 提供的专业代码编辑功能
    - 实时预览编辑效果
    - 代码格式化与高亮
3. **PDF 生成**：
    - 使用 wkhtmltopdf 保证高质量的 PDF 输出
    - 支持自定义页眉页脚
    - 多种页面尺寸选项
4. **配置管理**：
    - 前后端路由统一配置
    - 模板和 Schema 的版本管理

## 安装与配置

### 1. 安装扩展包

```bash
go get github.com/your-username/goravel-pdf
```

### 2. 配置

在 `config/app.go` 中添加配置：

```go
"pdf": map[string]interface{}{
    "prefix": "/api/pdf", // 前端API路由前缀
    "wkhtmltopdf_path": "/usr/local/bin/wkhtmltopdf", // wkhtmltopdf路径
    "default_template": "resources/views/default_template.html", // 默认模板路径
},
```

### 3. 注册服务提供者

在 `config/app.go` 的 `providers` 部分添加：

```go
"github.com/your-username/goravel-pdf/providers/PDFServiceProvider",
```

## 使用说明

### 1. 前端路由配置

后端路由会自动通过全局变量共享到前端：

```javascript
// 前端自动获取的路由
window.pdf_prefix = "/api/pdf";
window.getDefaultTemplate = "getDefaultTemplate";
window.saveTemplate = "saveTemplate/{id}";
window.saveHTML = "saveHTML/{id}";
window.generate = "generate/{id}";
```

### 2. 模板编辑

访问 `/pdf_design` 路由进入设计界面：

- **表单设计**：基于 JSON Schema 设计表单
- **模板编辑**：使用 Monaco Editor 编辑 HTML 模板
- **变量插入**：使用 `[[.variableName]]` 语法插入动态变量

### 3. API 接口

| 路由 | 方法 | 描述 |
|------|------|------|
| /api/pdf/getDefaultTemplate | GET | 获取默认模板和Schema |
| /api/pdf/saveTemplate/{id} | POST | 保存Schema定义 |
| /api/pdf/saveHTML/{id} | POST | 保存HTML模板 |
| /api/pdf/generate/{id} | GET | 生成PDF文件 |

## 开发指南

### 1. 扩展包开发

扩展包采用标准的 Goravel 扩展包结构：

```
goravel-pdf/
├── controllers/         # 控制器
├── providers/           # 服务提供者
├── resources/           # 资源文件
│   ├── views/           # 模板文件
│   └── assets/          # 静态资源
├── routes/              # 路由定义
└── config/              # 配置文件
```

### 2. 模板渲染

后端使用 Goravel 的模板引擎渲染 HTML：

```go
// 渲染web端页面
router.Get("/pdf_design", func(ctx http.Context) http.Response {
    // 读取HTML文件内容
    content := readHTMLFile("dashboard.html")
    
    return ctx.Response().View().Make("design.tmpl", map[string]any{
        "content": template.HTML(content),
    })
})
```

### 3. Monaco Editor 集成

前端使用 Monaco Editor 提供专业的代码编辑体验：

```javascript
// 初始化编辑器
editor = monaco.editor.create(document.getElementById('editorContainer'), {
    value: htmlContent,
    language: 'html',
    theme: 'vs-dark',
    automaticLayout: true,
    minimap: { enabled: true }
});
```

### 4. PDF 生成

使用 go-wkhtmltopdf 生成 PDF：

```go
func (p *PDFController) Generate(ctx http.Context) http.Response {
    // 获取HTML内容
    html := getTemplateHTML()
    
    // 初始化PDF生成器
    pdfg, err := wkhtmltopdf.NewPDFGenerator()
    if err != nil {
        return ctx.Response().Json(500, map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    // 添加HTML页面
    page := wkhtmltopdf.NewPage(html)
    pdfg.AddPage(page)
    
    // 生成PDF
    err = pdfg.Create()
    if err != nil {
        return ctx.Response().Json(500, map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    // 返回PDF文件
    return ctx.Response().Header("Content-Type", "application/pdf").
        Header("Content-Disposition", "attachment; filename=contract.pdf").
        Data(200, pdfg.Bytes())
}
```

## 未来计划

1. **JSON Schema 验证**：增强表单数据的验证功能
2. **模板版本控制**：实现模板的历史版本管理
3. **多主题支持**：为编辑器添加更多主题选项
4. **协作编辑**：支持多人同时编辑模板
5. **模板市场**：提供可共享的模板库

## 贡献指南

欢迎提交 Pull Request 或 Issue 报告问题。

## 许可证

MIT License