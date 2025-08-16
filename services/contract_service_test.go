package contract

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/goravel/framework/support/path"
	"github.com/hulutech-web/goravel_pdf_gen/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"html/template"
	"os"
	"testing"
	"time"
)

// TestContractHTMLToPDF 测试合同HTML转PDF功能
func TestContractHTMLToPDF(t *testing.T) {
	// 1. 准备合同模板（使用之前修正后的模板）
	contractTemplate := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>{{.ContractTitle}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            font-family: "SimHei", "Microsoft YaHei", sans-serif;
        }
        body {
            padding: 20px;
            line-height: 1.8;
            color: #333;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .title {
            text-align: center;
            margin-bottom: 60px;
            font-size: 24px;
            font-weight: bold;
        }
        /* 表格布局核心样式 */
		.contract-info {
			margin-bottom: 40px;
			font-size: 16px;
		}
		
		.info-table {
			width: 100%; /* 占满父容器宽度 */
			border-collapse: collapse; /* 合并边框，避免间隙 */
			table-layout: fixed; /* 固定布局，防止内容撑开表格 */
		}
		
		.info-tr {
			height: 45px; /* 固定行高，确保对齐 */
		}
		
		/* 标签单元格：固定宽度、加粗 */
		.info-label {
			width: 120px; /* 与原设计一致的固定宽度 */
			font-weight: bold;
			text-align: left; /* 左对齐 */
			padding: 0; /* 清除默认内边距 */
			vertical-align: bottom; /* 内容靠下，与下划线对齐 */
		}
		
		/* 内容单元格：下划线样式 */
		.info-content {
			padding: 0; /* 清除默认内边距 */
			vertical-align: bottom; /* 内容靠下，与下划线对齐 */
			border-bottom: 1px solid #999; /* 保留原有的下划线 */
			word-break: break-all; /* 内容过长时自动换行，不撑破表格 */
		}
        .terms {
            margin: 30px 0;
        }
        .term-title {
            font-weight: bold;
            margin: 20px 0 10px;
            font-size: 17px;
        }
        .term-content {
            text-indent: 2em;
            margin-bottom: 10px;
        }
		/* 表格核心样式：确保不换行、宽度适配 */
		.sign-table {
			width: 100%; /* 占满父容器宽度 */
			border-collapse: collapse; /* 合并边框（避免多余空隙） */
			table-layout: fixed; /* 固定表格布局，防止内容撑开导致换行 */
			margin-top: 80px; /* 保留原有的顶部间距 */
		}
		
		.sign-cell {
			width: 50%; /* 两列等宽 */
			padding: 0; /* 清除默认内边距 */
			vertical-align: top; /* 顶部对齐（避免上下错位） */
		}
		
		/* 保留原有的签名区域样式，无需修改 */
		.sign-party {
			text-align: start; /* 内容居中 */
			padding: 0 10px; /* 左右留少量间距，避免内容贴边 */
		}
		
		.sign-name {
			margin-top: 60px;
			border-bottom: 1px solid #999;
			padding-bottom: 5px;
			min-height: 30px; /* 确保签名线高度一致 */
		}
		/* 签名图片样式：限制最大尺寸，保持比例 */
		.sign-image {
			width: 180px; /* 最大宽度（根据签名区调整） */
			height: 80px; /* 最大高度（不超过容器） */
			object-fit: contain; /* 保持图片比例，不拉伸 */
		}
		.sign-image-container {
			margin-top: 30px; /* 与“签字/盖章”文字保持距离 */
			height: 80px; /* 固定签名区域高度 */
			width:120px;
		}
		/* 文字签名兜底：图片加载失败时显示 */
		.sign-name-fallback {
			margin-top: 60px;
			border-bottom: 1px solid #999;
			padding-bottom: 5px;
			min-height: 30px;
			display: none; /* 默认隐藏，图片加载失败时通过JS显示 */
		}
        .footer {
            margin-top: 100px;
            text-align: center;
            color: #666;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="title">{{.ContractTitle}}</div>

		<div class="contract-info">
			<table class="info-table">
				<!-- 第一行：合同编号 -->
				<tr class="info-tr">
					<td class="info-label">合同编号：</td>
					<td class="info-content">{{.ContractNo}}</td>
				</tr>
				<!-- 第二行：签订日期 -->
				<tr class="info-tr">
					<td class="info-label">签订日期：</td>
					<td class="info-content">{{.SignDate}}</td>
				</tr>
				<!-- 第三行：签订地点 -->
				<tr class="info-tr">
					<td class="info-label">签订地点：</td>
					<td class="info-content">{{.SignLocation}}</td>
				</tr>
			</table>
		</div>

        <div class="terms">
            <div class="term-title">甲方（{{.PartyARole}}）：
            <div class="term-content">名称：{{.PartyAName}}</div>
            <div class="term-content">身份证号/统一社会信用代码：{{.PartyAID}}</div>
            <div class="term-content">联系地址：{{.PartyAAddress}}</div>
            <div class="term-content">联系电话：{{.PartyAPhone}}</div>

            <div class="term-title">乙方（{{.PartyBRole}}）：</div>
            <div class="term-content">名称：{{.PartyBName}}</div>
            <div class="term-content">身份证号/统一社会信用代码：{{.PartyBID}}</div>
            <div class="term-content">联系地址：{{.PartyBAddress}}</div>
            <div class="term-content">联系电话：{{.PartyBPhone}}</div>
        </div>

        <div class="terms">
            <div class="term-title">第一条 合同标的</div>
            <div class="term-content">
                {{.Term1Content}}（租金每月人民币 {{.Amount}} 元整（大写：{{.AmountUpper}}））
            </div>

            <div class="term-title">第二条 权利与义务</div>
            <div class="term-content">
                <div>1. 甲方权利：</div><div>{{.PartyARights}}</div>
                <div>2. 甲方义务：</div><div>{{.PartyAObligations}}</div>
                <div>3. 乙方权利：</div><div>{{.PartyBRights}}</div>
                <div>4. 乙方义务：</div><div>{{.PartyBObligations}}</div>
            </div>

            <div class="term-title">第三条 违约责任</div>
            <div class="term-content">{{.LiabilityContent}}</div>

            <div class="term-title">第四条 争议解决</div>
            <div class="term-content">
                因本合同引起的或与本合同有关的任何争议，双方应首先通过友好协商解决；协商不成的，任何一方均有权向 {{.DisputeCourt}} 提起诉讼。
            </div>

            <div class="term-title">第五条 其他</div>
            <div class="term-content">
                <div>1. 本合同自双方签字（盖章）之日起生效。</div>
                <div>2. 本合同一式 {{.CopyCount}} 份，甲方执 {{.PartyACopy}} 份，乙方执 {{.PartyBCopy}} 份，具有同等法律效力。</div>
            </div>
        </div>

		
		<div class="sign-area">
			<table class="sign-table">
				<!-- 表格仅1行，2列分别放甲方、乙方 -->
				<tr>
					<td class="sign-cell">
						<div class="sign-party">
							<div>甲方（签字/盖章）：</div>
						  <!-- 电子签名图片：优先使用绝对路径或Base64编码 -->
							<div class="sign-image-container">
								<img src="{{.PartyASignImage}}" class="sign-image" alt="甲方签名">
							</div>
							<div class="sign-name">{{.PartyASign}}</div>
							<div>日期：{{.PartyASignDate}}</div>
						</div>
					</td>
					<td class="sign-cell">
						<div class="sign-party">
							<div>乙方（签字/盖章）：</div>
							<div class="sign-image-container">
									<img src="{{.PartyBSignImage}}" class="sign-image" alt="乙方签名">
							</div>
							<div class="sign-name">{{.PartyBSign}}</div>
							<div>日期：{{.PartyBSignDate}}</div>
						</div>
					</td>
				</tr>
			</table>
		</div>

        <div class="footer">
            注：本合同内容基于双方真实意思表示，电子版与纸质版具有同等法律效力。
        </div>
    </div>
</body>
</html>
`
	toBase64, err := TransImageToBase64("a.png")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"DEBUG - 转换图片失败": err,
		}).Error("DEBUG - 转换图片失败")
	}

	a_img := toBase64
	imageToBase64, err1 := TransImageToBase64("b.png")
	if err1 != nil {
		logrus.WithFields(logrus.Fields{
			"DEBUG - 转换图片失败": err1,
		}).Error("DEBUG - 转换图片失败")
	}

	b_img := imageToBase64
	// 2. 准备模板变量
	testVars := models.ContractVariables{
		ContractTitle:     "房屋租赁合同",
		ContractNo:        "ZL20240815001",
		SignDate:          "2024-08-15",
		SignLocation:      "北京市朝阳区建国路88号",
		PartyARole:        "出租方（房东）",
		PartyAName:        "张大山",
		PartyAID:          "1101051980XXXX1234",
		PartyAAddress:     "北京市海淀区中关村大街1号",
		PartyAPhone:       "13800138000",
		PartyBRole:        "承租方（租客）",
		PartyBName:        "李小四",
		PartyBID:          "3101041990XXXX5678",
		PartyBAddress:     "上海市浦东新区张江高科技园区",
		PartyBPhone:       "13900139000",
		Term1Content:      "甲方将其合法拥有的房屋（位于北京市朝阳区建国路88号，户型1室1厅）出租给乙方作为居住使用，租赁期限1年（2024-09-01至2025-08-31）",
		Amount:            "5000",
		AmountUpper:       "伍仟元整",
		PartyARights:      "1. 按时收取租金；\n2. 定期检查房屋使用情况（提前3天通知）；\n3. 租赁期满收回房屋\n",
		PartyAObligations: "1. 保证房屋产权清晰；\n2. 提供房屋合格证明；\n3. 房屋主体损坏7日内维修\n",
		PartyBRights:      "1. 合理使用房屋及设施；\n2. 要求甲方维修非人为损坏；\n3. 拒绝无故涨租\n",
		PartyBObligations: "1. 每月5日前付租金；\n2. 不擅自改结构/转租；\n3. 期满结清水电并保持房屋整洁\n",
		LiabilityContent:  "1. 甲方逾期交房，每日按月租1%付违约金；\n2. 乙方逾期付租，每日按欠缴1%付违约金（超15日甲方可解约）；\n3. 乙方转租需付2个月租金违约金\n",
		DisputeCourt:      "北京市朝阳区人民法院",
		CopyCount:         "2",
		PartyACopy:        "1",
		PartyBCopy:        "1",
		PartyASignImage:   template.URL(a_img),
		PartyBSignImage:   template.URL(b_img),
		PartyASign:        "张大山",
		PartyASignDate:    "2024-08-15",
		PartyBSign:        "李小四",
		PartyBSignDate:    "2024-08-15",
	}

	// 3. 渲染HTML（使用之前的模板渲染逻辑）
	tpl, err := template.New("contract").Parse(contractTemplate)
	assert.NoError(t, err, "模板解析失败: %v", err)

	var htmlBuf bytes.Buffer
	err = tpl.Execute(&htmlBuf, testVars)
	assert.NoError(t, err, "HTML渲染失败: %v", err)
	htmlContent := htmlBuf.String()

	//打印htmlContent
	//t.Logf("HTML内容:\n%s", htmlContent)
	// 4. 转换HTML为PDF
	pdfBytes, err := HTMLToPDF(htmlContent)
	assert.NoError(t, err, "HTML转PDF失败: %v", err)
	assert.NotEmpty(t, pdfBytes, "生成的PDF为空")

	// 5. （可选）保存PDF到本地（测试用）
	outputDir := "./test_pdfs"
	filename := fmt.Sprintf("contract_%s.pdf", time.Now().Format("20060102150405"))
	pdfPath, err := SavePDFToFile(pdfBytes, outputDir, filename)
	assert.NoError(t, err, "保存PDF失败: %v", err)
	t.Logf("PDF已保存至: %s", pdfPath) // 打印保存路径，方便查看测试结果
}

// TransImageToBase64签名图片转换为base64编码
func TransImageToBase64(imagePath string) (string, error) {
	imageBytes, err := os.ReadFile(path.Public(imagePath))
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageBytes), nil
}
