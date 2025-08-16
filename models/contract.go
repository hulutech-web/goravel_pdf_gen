package models

import (
	"bytes"
	"fmt"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
	html_template "html/template" // 导入模板包
	"text/template"
)

type Contract struct {
	orm.Model
	OrderID        uint            `gorm:"column:order_id;" form:"order_id" json:"order_id"` // 租赁合同关联订单
	LandlordID     uint            `gorm:"column:landlord_id;not null" form:"landlord_id" json:"landlord_id"`
	TenantID       uint            `gorm:"column:tenant_id;not null" form:"tenant_id" json:"tenant_id"` // 租赁合同需要租客
	Type           int8            `gorm:"column:type;not null" form:"type" json:"type"`                // 1-代理合同 2-租赁合同
	Content        string          `gorm:"type:text;not null" form:"content" json:"content"`
	TenantSign     string          `gorm:"column:tenant_sign" form:"tenant_sign" json:"tenant_sign"`
	LandlordSign   string          `gorm:"column:landlord_sign" form:"landlord_sign" json:"landlord_sign"`
	SignedAt       carbon.DateTime `gorm:"column:signed_at" form:"signed_at" json:"signed_at"`
	ExpiredAt      carbon.DateTime `gorm:"expired_at" from:"expired_at" json:"expired_at"`
	SignedLocation string          `gorm:"column:signed_location" form:"signed_location" json:"signed_location"`
	PaperContract  string          `gorm:"column:paper_contract" form:"paper_contract" json:"paper_contract"`
	// 关系
}

// 合同模板变量结构体（与HTML模板中的占位符对应）
type ContractVariables struct {
	ContractTitle     string            `json:"contract_title"`      // 合同标题
	ContractNo        string            `json:"contract_no"`         // 合同编号
	SignDate          string            `json:"sign_date"`           // 签订日期
	SignLocation      string            `json:"sign_location"`       // 签订地点
	PartyARole        string            `json:"party_a_role"`        // 甲方角色
	PartyAName        string            `json:"party_a_name"`        // 甲方名称
	PartyAID          string            `json:"party_a_id"`          // 甲方身份证
	PartyAAddress     string            `json:"party_a_address"`     // 甲方地址
	PartyAPhone       string            `json:"party_a_phone"`       // 甲方电话
	PartyBRole        string            `json:"party_b_role"`        // 乙方角色
	PartyBName        string            `json:"party_b_name"`        // 乙方名称
	PartyBID          string            `json:"party_b_id"`          // 乙方身份证
	PartyBAddress     string            `json:"party_b_address"`     // 乙方地址
	PartyBPhone       string            `json:"party_b_phone"`       // 乙方电话
	Term1Content      string            `json:"term_1_content"`      // 第一条内容
	Amount            string            `json:"amount"`              // 金额
	AmountUpper       string            `json:"amount_upper"`        // 金额大写
	PartyARights      string            `json:"party_a_rights"`      // 甲方权利
	PartyAObligations string            `json:"party_a_obligations"` // 甲方义务
	PartyBRights      string            `json:"party_b_rights"`      // 乙方权利
	PartyBObligations string            `json:"party_b_obligations"` // 乙方义务
	LiabilityContent  string            `json:"liability_content"`   // 违约责任
	DisputeCourt      string            `json:"dispute_court"`       // 争议管辖法院
	CopyCount         string            `json:"copy_count"`          // 合同份数
	PartyACopy        string            `json:"party_a_copy"`        // 甲方持有份数
	PartyBCopy        string            `json:"party_b_copy"`        // 乙方持有份数
	PartyASign        string            `json:"party_a_sign"`        // 甲方签名
	PartyASignImage   html_template.URL `json:"party_a_sign_image"`  // 甲方签名图片
	PartyBSignImage   html_template.URL `json:"party_b_sign_image"`  // 乙方签名图片
	PartyBSign        string            `json:"party_b_sign"`        // 乙方签名
	PartyASignDate    string            `json:"party_a_sign_date"`   // 甲方签名日期
	PartyBSignDate    string            `json:"party_b_sign_date"`   // 乙方签名日期
}

// 生成合同编号（规则：类型+日期+ID，如 ZL20240520001）
func (c *Contract) GenerateContractNo() string {
	prefix := "DL" // 代理合同前缀
	if c.Type == 2 {
		prefix = "ZL" // 租赁合同前缀
	}
	date := carbon.Now().Format("Ymd")
	return prefix + date + fmt.Sprintf("%03d", c.ID)
}

// 渲染合同模板（将变量替换到HTML模板中）
// 【关键修改】渲染合同模板：手动构建变量Map，匹配模板占位符
func (c *Contract) RenderTemplate(templateContent string, variables ContractVariables) (string, error) {
	// 手动构建“模板占位符名→结构体字段值”的映射Map
	variablesMap := map[string]interface{}{
		"contract_title":      variables.ContractTitle,
		"contract_no":         variables.ContractNo,
		"sign_date":           variables.SignDate,
		"sign_location":       variables.SignLocation,
		"party_a_role":        variables.PartyARole,
		"party_a_name":        variables.PartyAName,
		"party_a_id":          variables.PartyAID,
		"party_a_address":     variables.PartyAAddress,
		"party_a_phone":       variables.PartyAPhone,
		"party_b_role":        variables.PartyBRole,
		"party_b_name":        variables.PartyBName,
		"party_b_id":          variables.PartyBID,
		"party_b_address":     variables.PartyBAddress,
		"party_b_phone":       variables.PartyBPhone,
		"term_1_content":      variables.Term1Content,
		"amount":              variables.Amount,
		"amount_upper":        variables.AmountUpper,
		"party_a_rights":      variables.PartyARights,
		"party_a_obligations": variables.PartyAObligations,
		"party_b_rights":      variables.PartyBRights,
		"party_b_obligations": variables.PartyBObligations,
		"liability_content":   variables.LiabilityContent,
		"dispute_court":       variables.DisputeCourt,
		"copy_count":          variables.CopyCount,
		"party_a_copy":        variables.PartyACopy,
		"party_b_copy":        variables.PartyBCopy,
		"party_a_sign":        variables.PartyASign,
		"party_b_sign":        variables.PartyBSign,
		"party_a_sign_date":   variables.PartyASignDate,
		"party_b_sign_date":   variables.PartyBSignDate,
	}

	// 解析模板并传入映射Map渲染
	tpl, err := template.New("contract").Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	// 此处传入的是手动构建的Map，模板可直接通过{{占位符名}}访问
	if err := tpl.Execute(&buf, variablesMap); err != nil {
		return "", err
	}

	return buf.String(), nil
}
