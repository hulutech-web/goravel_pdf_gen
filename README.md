<p align="center">
  <img src="https://github.com/hulutech-web/goravel_pdf_gen/blob/master/images/logo.png?raw=true" width="300"  />
</p>

# goravel_pdf_gen
#### A simple package to generate PDFs in Go using the goravel framework
#### 配置内容
- 定义表单字段，生成migration文件，并创建model，创建controller，创建route
- 字段包含这个3个，id,name,params,其中params为一个json结构，它定义了pdf动态的数据内容
- 举例说明，合同场景，合同标题，创建时间，合同内容，合同编号，合同金额，合同期限，合同甲方，合同乙方，合同丙方，合同丁方，合同戊方，合同己方，合同庚方，合同辛方，合同壬方，合同癸方，合同十一方，合同十二方，合同十三方，合同十四方，合同十五方，合同十六方，合同十七方，
- 模型映射，使用反射，将现有模型的字段，映射到json中的字段，通过修改config/goravel_pdf_gen.go进行map绑定，实现转换
- 路由，1、定义表，2、定义字段