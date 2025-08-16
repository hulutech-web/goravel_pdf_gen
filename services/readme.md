### 电子合同
#### 三方库
- 思路：采用html动态生成pdf，其他思路兼容性差，无法满足要求   
``github.com/SebastiaanKlippert/go-wkhtmltopdf``
  1. [go-wkhtmltopdf](https://github.com/SebastiaanKlippert/go-wkhtmltopdf.git) 需要安装本地环境，根据操作系统安装，mac linux windows都需要本地安装
  2. API方式调用
  3. 程序配置
     - 配置
       ```go

        func HTMLToPDF(htmlContent string) ([]byte, error) {
        pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
        if err != nil {
        return nil, fmt.Errorf("初始化PDF生成器失败: %w", err)
        }
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
              ```
     - 注意事项  
     图片问题，需要使用base64，同时因golang 模板限制，html/template 对不安全内容的转义占位符，``#ZgotmplZ``。    
     当直接传递 []byte 或未标记为安全的二进制数据时，模板引擎会阻止潜在的安全风险，导致 Base64 被替换。   
     解决方案是：定义如下，告诉golang，明确告知模板引擎该内容是安全的 URL 协议，禁止转义
     ```go
      PartyASignImage   html_template.URL `json:"party_a_sign_image"`  // 甲方签名图片
      PartyBSignImage   html_template.URL `json:"party_b_sign_image"`  // 乙方签名图片
    ```