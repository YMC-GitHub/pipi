# pipi - 纯函数图像处理与定位库

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ymc-github/pipi)](https://goreportcard.com/report/github.com/ymc-github/pipi)

pipi 是一个专注于文本水印定位和图像处理的 Go 语言纯函数库。它提供了类型安全、无副作用的函数，用于处理颜色、边界框、字体大小、文本对齐等常见图像处理任务。

## ✨ 特性

- 🎨 **颜色处理** - 支持颜色名、十六进制、RGBA 格式解析
- 📐 **边界框管理** - 支持 XYXY、XYWH、CXCY 多种坐标格式
- 📏 **字体大小计算** - 支持 px、%、rem、bh 单位
- 🎯 **文本对齐** - 灵活的对齐方式，支持偏移量
- 🖌️ **描边效果** - 自动生成描边偏移量
- ✅ **类型安全** - 纯函数设计，无副作用
- 🧪 **测试覆盖** - 完整的单元测试和基准测试
- 📦 **零依赖** - 仅使用 Go 标准库

## 📦 安装

```bash
go get github.com/ymc-github/pipi
```

## 🚀 快速开始

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pipi"
)

func main() {
    // 解析颜色
    color := pipi.ParseColor("#FF0000")
    fmt.Printf("Red color: %+v\n", color)
    
    // 解析边界框
    x1, y1, x2, y2, _ := pipi.ParseBBox("10,20,100,200", pipi.BBoxXYXY)
    fmt.Printf("BBox: (%.0f, %.0f) to (%.0f, %.0f)\n", x1, y1, x2, y2)
    
    // 计算字体大小
    fontSize := pipi.CalculateFontSize(2, pipi.UnitREM, 16, 100)
    fmt.Printf("Font size: %.0fpx\n", fontSize)
    
    // 解析对齐方式
    align, _ := pipi.ParseAlignment("bottom+10px", 16)
    fmt.Printf("Alignment: %+v\n", align)
}
```

## 📚 API 文档

### 颜色处理

#### ParseColor
解析颜色字符串，支持多种格式。

```go
func ParseColor(colorStr string) color.Color
```

**支持格式：**
- 颜色名：`black`, `white`, `red`, `green`, `blue`, `yellow`, `cyan`, `magenta`, `transparent`
- 十六进制：`#FF0000`, `#FF0000FF`
- RGBA：`rgba(255,0,0,128)`

**示例：**
```go
c1 := pipi.ParseColor("red")           // RGB(255,0,0)
c2 := pipi.ParseColor("#00FF00")       // RGB(0,255,0)
c3 := pipi.ParseColor("rgba(0,0,255,128)") // RGBA(0,0,255,128)
```

#### ParseColorWithDefault
解析颜色，失败时返回默认颜色。

```go
func ParseColorWithDefault(colorStr string, defaultColor color.Color) color.Color
```

#### ColorToRGBA
将 color.Color 转换为 RGBA 分量。

```go
func ColorToRGBA(c color.Color) (r, g, b, a uint8)
```

### 边界框处理

#### BBoxType 类型
```go
const (
    BBoxXYXY pipi.BBoxType = "xyxy" // x1,y1,x2,y2
    BBoxXYWH pipi.BBoxType = "xywh" // x,y,width,height
    BBoxCXCY pipi.BBoxType = "cxcy" // center-x,center-y,width,height
)
```

#### ParseBBox
解析边界框字符串。

```go
func ParseBBox(bboxStr string, bboxType BBoxType) (x1, y1, x2, y2 float64, err error)
```

**示例：**
```go
// XYXY 格式
x1, y1, x2, y2, _ := pipi.ParseBBox("10,20,100,200", pipi.BBoxXYXY)

// XYWH 格式（自动转换为 XYXY）
x1, y1, x2, y2, _ := pipi.ParseBBox("10,20,90,180", pipi.BBoxXYWH)

// CXCY 格式（自动转换为 XYXY）
x1, y1, x2, y2, _ := pipi.ParseBBox("55,110,90,180", pipi.BBoxCXCY)
```

#### ConvertBBox
在不同边界框格式间转换。

```go
func ConvertBBox(x1, y1, x2, y2 float64, toType BBoxType) [4]float64
```

#### ValidateBBox
验证边界框是否有效。

```go
func ValidateBBox(x1, y1, x2, y2 float64) bool
```

#### GetDefaultBBox
获取默认边界框（整个图像）。

```go
func GetDefaultBBox(imgWidth, imgHeight float64) BBox
```

### 字体大小计算

#### FontSizeUnit 类型
```go
const (
    UnitPX  pipi.FontSizeUnit = "px"  // 绝对像素
    UnitPct pipi.FontSizeUnit = "%"   // 边界框高度百分比
    UnitREM pipi.FontSizeUnit = "rem" // 相对根字体大小
    UnitBH  pipi.FontSizeUnit = "bh"  // 边界框高度倍数
)
```

#### CalculateFontSize
根据单位和上下文计算实际字体大小。

```go
func CalculateFontSize(sizeValue float64, unit FontSizeUnit, rootSize, bboxHeight float64) float64
```

**示例：**
```go
// px：绝对像素
size := pipi.CalculateFontSize(24, pipi.UnitPX, 16, 100) // 返回 24

// rem：相对根字体（16px * 2 = 32px）
size := pipi.CalculateFontSize(2, pipi.UnitREM, 16, 100) // 返回 32

// %：边界框高度的百分比（100px * 50% = 50px）
size := pipi.CalculateFontSize(50, pipi.UnitPct, 16, 100) // 返回 50

// bh：边界框高度倍数（100px * 0.5 = 50px）
size := pipi.CalculateFontSize(0.5, pipi.UnitBH, 16, 100) // 返回 50
```

#### ParseAndCalculateFontSize
解析大小字符串并计算实际字体大小。

```go
func ParseAndCalculateFontSize(sizeStr string, rootSize, bboxHeight float64) (fontSize float64, unit FontSizeUnit, err error)
```

### 文本对齐

#### Alignment 结构
```go
type Alignment struct {
    Type      string  // 对齐类型：baseline, top, bottom, middle, left, right, center
    Offset    float64 // 像素偏移
    OffsetRel float64 // 相对偏移（0-1）
    IsRel     bool    // 是否使用相对偏移
}
```

#### ParseAlignment
解析对齐字符串。

```go
func ParseAlignment(alignStr string, rootSize float64) (Alignment, error)
```

**支持格式：**
- `baseline` - 基线对齐
- `top` - 顶部对齐
- `bottom` - 底部对齐
- `middle` / `center` - 居中对齐
- `left` - 左对齐
- `right` - 右对齐
- 支持偏移：`bottom+10px`, `middle-5%`, `top+2rem`

**示例：**
```go
// 底部对齐，向上偏移 10px
align, _ := pipi.ParseAlignment("bottom-10px", 16)

// 居中对齐，向右偏移 5%
align, _ := pipi.ParseAlignment("center+5%", 16)

// 顶部对齐，向下偏移 2rem
align, _ := pipi.ParseAlignment("top+2rem", 16)
```

#### CalculateAlignedPositionX
计算水平对齐的 X 坐标。

```go
func CalculateAlignedPositionX(bboxMinX, bboxMaxX, textWidth float64, align Alignment) float64
```

#### CalculateAlignedPositionY
计算垂直对齐的 Y 坐标。

```go
func CalculateAlignedPositionY(bboxMinY, bboxMaxY, textHeight float64, align Alignment) float64
```

#### CalculateAlignedPosition
计算完整的对齐位置。

```go
func CalculateAlignedPosition(bbox BBox, textWidth, textHeight float64, alignX, alignY Alignment) (x, y float64)
```

### 文本边界检查

#### CheckTextBounds
检查文本是否在边界框内。

```go
func CheckTextBounds(x, y, textWidth, textHeight float64, bbox BBox) bool
```

#### GetTextBounds
获取文本的边界框。

```go
func GetTextBounds(x, y, textWidth, textHeight float64) BBox
```

### 描边处理

#### GenerateStrokeOffsets
生成描边偏移量向量。

```go
func GenerateStrokeOffsets(strokeWidth float64) []struct{ Dx, Dy float64 }
```

**示例：**
```go
offsets := pipi.GenerateStrokeOffsets(2.0)
for _, offset := range offsets {
    // 在偏移位置绘制文本
    drawText(x + offset.Dx, y + offset.Dy)
}
```

#### GetStrokeColor
获取描边颜色（如果描边宽度为 0 或描边颜色为 nil，则返回文本颜色）。

```go
func GetStrokeColor(strokeWidth float64, strokeColor, textColor color.Color) color.Color
```

### 辅助函数

#### ParseSizeValue
解析大小字符串（如 "24px", "2rem", "10%"）。

```go
func ParseSizeValue(sizeStr string) (float64, string, error)
```

#### ParseImageSize
解析图片尺寸字符串（如 "1280x720"）。

```go
func ParseImageSize(sizeStr string) (width, height int, err error)
```

#### IsValidSizeUnit
检查大小单位是否有效。

```go
func IsValidSizeUnit(unit string) bool
```

#### NormalizeCommand
规范化命令字符串（将 "center" 转换为 "middle"）。

```go
func NormalizeCommand(cmd string) string
```

#### ParseCoordinate
解析可能带单位的坐标值。

```go
func ParseCoordinate(coordStr string, rootSize float64) (float64, error)
```

## 💡 完整示例

### 文本水印定位

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pipi"
)

func main() {
    // 图像尺寸
    imgWidth, imgHeight := 1920.0, 1080.0
    
    // 边界框（右下角区域）
    bbox := pipi.BBox{1500, 800, 1900, 1000}
    
    // 文本尺寸（模拟）
    textWidth, textHeight := 300.0, 50.0
    
    // 对齐方式：右下角，向外偏移 20px
    alignX, _ := pipi.ParseAlignment("right-20px", 16)
    alignY, _ := pipi.ParseAlignment("bottom-20px", 16)
    
    // 计算位置
    x, y := pipi.CalculateAlignedPosition(bbox, textWidth, textHeight, alignX, alignY)
    
    fmt.Printf("Text position: (%.0f, %.0f)\n", x, y)
    
    // 检查是否在边界框内
    inside := pipi.CheckTextBounds(x, y, textWidth, textHeight, bbox)
    fmt.Printf("Text inside bbox: %v\n", inside)
}
```

### 字体大小自适应

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pipi"
)

func main() {
    bboxHeight := 200.0  // 边界框高度 200px
    rootSize := 16.0      // 根字体大小 16px
    
    // 使用 rem 单位
    fontSize, unit, _ := pipi.ParseAndCalculateFontSize("2rem", rootSize, bboxHeight)
    fmt.Printf("2rem = %.0fpx (%s)\n", fontSize, unit)
    
    // 使用百分比
    fontSize, unit, _ = pipi.ParseAndCalculateFontSize("50%", rootSize, bboxHeight)
    fmt.Printf("50%% = %.0fpx (%s)\n", fontSize, unit)
    
    // 使用 bh 单位
    fontSize, unit, _ = pipi.ParseAndCalculateFontSize("0.5bh", rootSize, bboxHeight)
    fmt.Printf("0.5bh = %.0fpx (%s)\n", fontSize, unit)
}
```

### 颜色处理

```go
package main

import (
    "fmt"
    "image/color"
    "github.com/ymc-github/pipi"
)

func main() {
    colors := []string{
        "red",
        "#00FF00",
        "rgba(0,0,255,128)",
        "unknown",
    }
    
    for _, c := range colors {
        parsed := pipi.ParseColor(c)
        r, g, b, a := pipi.ColorToRGBA(parsed)
        fmt.Printf("%-20s -> RGBA(%3d,%3d,%3d,%3d)\n", c, r, g, b, a)
    }
}
```

## 🧪 测试

运行测试套件：

```bash
# 运行所有测试
go test -v

# 运行基准测试
go test -bench=. -benchmem

# 生成覆盖率报告
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```
## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

## 📧 联系方式

- 项目主页: [https://github.com/ymc-github/pipi](https://github.com/ymc-github/pipi)
- 问题反馈: [Issues](https://github.com/ymc-github/pipi/issues)
