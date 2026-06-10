# pipi - Pure Function Image Processing & Positioning Library

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ymc-github/pipi)](https://goreportcard.com/report/github/ymc-github/pipi)

pipi is a Go pure function library focused on text watermark positioning and image processing. It provides type-safe, side-effect-free functions for common image processing tasks such as color handling, bounding boxes, font size calculation, and text alignment.

## ✨ Features

- 🎨 **Color Processing** - Supports color names, hex, and RGBA format parsing
- 📐 **Bounding Box Management** - Supports XYXY, XYWH, CXCY coordinate formats
- 📏 **Font Size Calculation** - Supports px, %, rem, bh units
- 🎯 **Text Alignment** - Flexible alignment with offset support
- 🖌️ **Stroke Effect** - Automatic stroke offset generation
- ✅ **Type Safety** - Pure function design with no side effects
- 🧪 **Test Coverage** - Complete unit tests and benchmarks
- 📦 **Zero Dependencies** - Uses only Go standard library

## 📦 Installation

```bash
go get github.com/ymc-github/pipi
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pipi"
)

func main() {
    // Parse color
    color := pipi.ParseColor("#FF0000")
    fmt.Printf("Red color: %+v\n", color)
    
    // Parse bounding box
    x1, y1, x2, y2, _ := pipi.ParseBBox("10,20,100,200", pipi.BBoxXYXY)
    fmt.Printf("BBox: (%.0f, %.0f) to (%.0f, %.0f)\n", x1, y1, x2, y2)
    
    // Calculate font size
    fontSize := pipi.CalculateFontSize(2, pipi.UnitREM, 16, 100)
    fmt.Printf("Font size: %.0fpx\n", fontSize)
    
    // Parse alignment
    align, _ := pipi.ParseAlignment("bottom+10px", 16)
    fmt.Printf("Alignment: %+v\n", align)
}
```

## 📚 API Documentation

### Color Processing

#### ParseColor
Parses color string, supporting multiple formats.

```go
func ParseColor(colorStr string) color.Color
```

**Supported formats:**
- Color names: `black`, `white`, `red`, `green`, `blue`, `yellow`, `cyan`, `magenta`, `transparent`
- Hexadecimal: `#FF0000`, `#FF0000FF`
- RGBA: `rgba(255,0,0,128)`

**Example:**
```go
c1 := pipi.ParseColor("red")                 // RGB(255,0,0)
c2 := pipi.ParseColor("#00FF00")             // RGB(0,255,0)
c3 := pipi.ParseColor("rgba(0,0,255,128)")   // RGBA(0,0,255,128)
```

#### ParseColorWithDefault
Parses color string, returns default color on failure.

```go
func ParseColorWithDefault(colorStr string, defaultColor color.Color) color.Color
```

#### ColorToRGBA
Converts color.Color to RGBA components.

```go
func ColorToRGBA(c color.Color) (r, g, b, a uint8)
```

### Bounding Box Processing

#### BBoxType Constants
```go
const (
    BBoxXYXY pipi.BBoxType = "xyxy" // x1,y1,x2,y2
    BBoxXYWH pipi.BBoxType = "xywh" // x,y,width,height
    BBoxCXCY pipi.BBoxType = "cxcy" // center-x,center-y,width,height
)
```

#### ParseBBox
Parses bounding box string.

```go
func ParseBBox(bboxStr string, bboxType BBoxType) (x1, y1, x2, y2 float64, err error)
```

**Example:**
```go
// XYXY format
x1, y1, x2, y2, _ := pipi.ParseBBox("10,20,100,200", pipi.BBoxXYXY)

// XYWH format (automatically converts to XYXY)
x1, y1, x2, y2, _ := pipi.ParseBBox("10,20,90,180", pipi.BBoxXYWH)

// CXCY format (automatically converts to XYXY)
x1, y1, x2, y2, _ := pipi.ParseBBox("55,110,90,180", pipi.BBoxCXCY)
```

#### ConvertBBox
Converts bounding box between different formats.

```go
func ConvertBBox(x1, y1, x2, y2 float64, toType BBoxType) [4]float64
```

#### ValidateBBox
Validates if bounding box is valid.

```go
func ValidateBBox(x1, y1, x2, y2 float64) bool
```

#### GetDefaultBBox
Returns default bounding box (entire image).

```go
func GetDefaultBBox(imgWidth, imgHeight float64) BBox
```

### Font Size Calculation

#### FontSizeUnit Constants
```go
const (
    UnitPX  pipi.FontSizeUnit = "px"  // absolute pixels
    UnitPct pipi.FontSizeUnit = "%"   // percentage of bbox height
    UnitREM pipi.FontSizeUnit = "rem" // relative to root size
    UnitBH  pipi.FontSizeUnit = "bh"  // box height multiplier
)
```

#### CalculateFontSize
Calculates actual font size based on unit and context.

```go
func CalculateFontSize(sizeValue float64, unit FontSizeUnit, rootSize, bboxHeight float64) float64
```

**Example:**
```go
// px: absolute pixels
size := pipi.CalculateFontSize(24, pipi.UnitPX, 16, 100) // returns 24

// rem: relative to root font (16px * 2 = 32px)
size := pipi.CalculateFontSize(2, pipi.UnitREM, 16, 100) // returns 32

// %: percentage of bbox height (100px * 50% = 50px)
size := pipi.CalculateFontSize(50, pipi.UnitPct, 16, 100) // returns 50

// bh: box height multiplier (100px * 0.5 = 50px)
size := pipi.CalculateFontSize(0.5, pipi.UnitBH, 16, 100) // returns 50
```

#### ParseAndCalculateFontSize
Parses size string and calculates actual font size.

```go
func ParseAndCalculateFontSize(sizeStr string, rootSize, bboxHeight float64) (fontSize float64, unit FontSizeUnit, err error)
```

### Text Alignment

#### Alignment Struct
```go
type Alignment struct {
    Type      string  // alignment type: baseline, top, bottom, middle, left, right, center
    Offset    float64 // pixel offset
    OffsetRel float64 // relative offset (0-1)
    IsRel     bool    // whether using relative offset
}
```

#### ParseAlignment
Parses alignment string.

```go
func ParseAlignment(alignStr string, rootSize float64) (Alignment, error)
```

**Supported formats:**
- `baseline` - baseline alignment
- `top` - top alignment
- `bottom` - bottom alignment
- `middle` / `center` - center alignment
- `left` - left alignment
- `right` - right alignment
- Supports offsets: `bottom+10px`, `middle-5%`, `top+2rem`

**Example:**
```go
// Bottom alignment, offset up by 10px
align, _ := pipi.ParseAlignment("bottom-10px", 16)

// Center alignment, offset right by 5%
align, _ := pipi.ParseAlignment("center+5%", 16)

// Top alignment, offset down by 2rem
align, _ := pipi.ParseAlignment("top+2rem", 16)
```

#### CalculateAlignedPositionX
Calculates X coordinate for horizontal alignment.

```go
func CalculateAlignedPositionX(bboxMinX, bboxMaxX, textWidth float64, align Alignment) float64
```

#### CalculateAlignedPositionY
Calculates Y coordinate for vertical alignment.

```go
func CalculateAlignedPositionY(bboxMinY, bboxMaxY, textHeight float64, align Alignment) float64
```

#### CalculateAlignedPosition
Calculates complete aligned position.

```go
func CalculateAlignedPosition(bbox BBox, textWidth, textHeight float64, alignX, alignY Alignment) (x, y float64)
```

### Text Bounds Checking

#### CheckTextBounds
Checks if text is within bounding box.

```go
func CheckTextBounds(x, y, textWidth, textHeight float64, bbox BBox) bool
```

#### GetTextBounds
Returns the bounding box of rendered text.

```go
func GetTextBounds(x, y, textWidth, textHeight float64) BBox
```

### Stroke Processing

#### GenerateStrokeOffsets
Generates stroke offset vectors.

```go
func GenerateStrokeOffsets(strokeWidth float64) []struct{ Dx, Dy float64 }
```

**Example:**
```go
offsets := pipi.GenerateStrokeOffsets(2.0)
for _, offset := range offsets {
    // Draw text at offset position
    drawText(x + offset.Dx, y + offset.Dy)
}
```

#### GetStrokeColor
Returns stroke color (returns text color if stroke width is 0 or stroke color is nil).

```go
func GetStrokeColor(strokeWidth float64, strokeColor, textColor color.Color) color.Color
```

### Helper Functions

#### ParseSizeValue
Parses size string (e.g., "24px", "2rem", "10%").

```go
func ParseSizeValue(sizeStr string) (float64, string, error)
```

#### ParseImageSize
Parses image size string (e.g., "1280x720").

```go
func ParseImageSize(sizeStr string) (width, height int, err error)
```

#### IsValidSizeUnit
Checks if size unit is valid.

```go
func IsValidSizeUnit(unit string) bool
```

#### NormalizeCommand
Normalizes command string (converts "center" to "middle").

```go
func NormalizeCommand(cmd string) string
```

#### ParseCoordinate
Parses coordinate value that may have unit suffixes.

```go
func ParseCoordinate(coordStr string, rootSize float64) (float64, error)
```

## 💡 Complete Examples

### Text Watermark Positioning

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pipi"
)

func main() {
    // Image dimensions
    imgWidth, imgHeight := 1920.0, 1080.0
    
    // Bounding box (bottom-right area)
    bbox := pipi.BBox{1500, 800, 1900, 1000}
    
    // Text dimensions (simulated)
    textWidth, textHeight := 300.0, 50.0
    
    // Alignment: bottom-right, offset outward by 20px
    alignX, _ := pipi.ParseAlignment("right-20px", 16)
    alignY, _ := pipi.ParseAlignment("bottom-20px", 16)
    
    // Calculate position
    x, y := pipi.CalculateAlignedPosition(bbox, textWidth, textHeight, alignX, alignY)
    
    fmt.Printf("Text position: (%.0f, %.0f)\n", x, y)
    
    // Check if inside bounding box
    inside := pipi.CheckTextBounds(x, y, textWidth, textHeight, bbox)
    fmt.Printf("Text inside bbox: %v\n", inside)
}
```

### Adaptive Font Size

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pipi"
)

func main() {
    bboxHeight := 200.0  // bbox height 200px
    rootSize := 16.0      // root font size 16px
    
    // Using rem unit
    fontSize, unit, _ := pipi.ParseAndCalculateFontSize("2rem", rootSize, bboxHeight)
    fmt.Printf("2rem = %.0fpx (%s)\n", fontSize, unit)
    
    // Using percentage
    fontSize, unit, _ = pipi.ParseAndCalculateFontSize("50%", rootSize, bboxHeight)
    fmt.Printf("50%% = %.0fpx (%s)\n", fontSize, unit)
    
    // Using bh unit
    fontSize, unit, _ = pipi.ParseAndCalculateFontSize("0.5bh", rootSize, bboxHeight)
    fmt.Printf("0.5bh = %.0fpx (%s)\n", fontSize, unit)
}
```

### Color Processing

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

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Generate coverage report
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 🙏 Acknowledgments

Thanks to all developers who have contributed to this project!

## 📧 Contact

- Project Home: [https://github.com/ymc-github/pipi](https://github.com/ymc-github/pipi)
- Issue Tracker: [Issues](https://github.com/ymc-github/pipi/issues)
