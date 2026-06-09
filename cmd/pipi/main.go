package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	Version         = "v1.4.0"
	DefaultRootSize = 16.0
	Usage           = `
Pipi - Watermark Text Tool

Usage:
  pipi text [options]
  pipi help
  pipi version

Commands:
  text                Draw watermark text

Options for 'text' command:
  --content TEXT      Text content (required)
  --size SIZE         Font size, support px, %, rem, bh
                      px: absolute pixels (default: 24px)
                      %: relative to bbox height (percentage)
                      rem: relative to root-size (CSS standard, default: 16px)
                      bh: relative to bbox height (box height multiplier)
                      Examples: 24px, 5%, 2rem, 0.1bh
  --root-size PX      Root font size for rem unit (default: 16px)
  --color COLOR       Text color, support color names or hex
                      Examples: red, #FF0000, rgba(255,0,0,255) (default: white)
  --stroke WIDTH      Text stroke width (default: 0, no stroke)
  --stroke-width WIDTH Alias for --stroke
  --stroke-color COLOR Stroke color (default: same as text color)
  --bbox COORDS       Bounding box coordinates
  --bbox-type TYPE    Bbox type: xyxy, xywh, cxcy (default: xyxy)
  --align-y ALIGN     Y-axis alignment: baseline, top, bottom, middle
                      Support px, %, rem offset (e.g., bottom+10px, middle-5%)
  --align-x ALIGN     X-axis alignment: left, right, middle, baseline
                      Support px, %, rem offset (e.g., left+10px, middle-5%)
  --font PATH         Font file path or system font name
  --input PATH        Input image path (creates transparent if empty)
  --imgz SIZE         Image size when --input is empty (default: 1280x720)
                      Format: WIDTHxHEIGHT
  --output PATH       Output image path (required)
  --debug             Enable debug output

Examples:
  # CSS standard rem
  pipi text --content "Hello" --size 2rem --output out.png
  
  # Box height multiplier
  pipi text --content "Hello" --size 0.1bh --output out.png
  
  # TTC font support (e.g., Microsoft YaHei)
  pipi text --content "微软雅黑测试" --font /fonts/msyh.ttc --output out.png
  
  # Stroke with different color
  pipi text --content "STROKE" --color white --stroke 3 --stroke-color black --output out.png
`
)

// 尺寸定义
type Size struct {
	Width, Height int
}

// 边界框类型
type BBoxType string

const (
	BBoxXYXY BBoxType = "xyxy"
	BBoxXYWH BBoxType = "xywh"
	BBoxCXCY BBoxType = "cxcy"
)

// 对齐方式
type Alignment struct {
	Type      string  // baseline, top, bottom, middle, left, right
	Offset    float64 // 偏移值（像素）
	OffsetRel float64 // 相对偏移（百分比，0-1）
	IsRel     bool    // 是否使用相对偏移
}

// 配置结构
type TextConfig struct {
	Content      string
	Size         float64
	SizeUnit     string  // px, %, rem, bh
	RootSize     float64 // rem 基准大小
	Color        color.Color
	StrokeWidth  float64
	StrokeColor  color.Color
	BBox         [4]float64
	BBoxType     BBoxType
	AlignY       Alignment
	AlignX       Alignment
	FontPath     string
	InputPath    string
	ImageSize    Size
	OutputPath   string
	Debug        bool
}

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "help":
			fmt.Print(Usage)
			return
		case "version":
			fmt.Println(Version)
			return
		}
	}

	if len(os.Args) < 2 || os.Args[1] != "text" {
		fmt.Println("Error: Invalid command. Use 'pipi text', 'pipi help', or 'pipi version'")
		fmt.Print(Usage)
		os.Exit(1)
	}

	// 解析 text 命令参数
	textCmd := flag.NewFlagSet("text", flag.ExitOnError)

	content := textCmd.String("content", "", "Text content")
	size := textCmd.String("size", "24px", "Font size")
	rootSize := textCmd.Float64("root-size", DefaultRootSize, "Root font size for rem unit")
	colorStr := textCmd.String("color", "white", "Text color")
	stroke := textCmd.String("stroke", "0", "Stroke width")
	strokeWidth := textCmd.String("stroke-width", "", "Stroke width (alias)")
	strokeColorStr := textCmd.String("stroke-color", "", "Stroke color (default: same as text color)")
	bbox := textCmd.String("bbox", "", "Bounding box")
	bboxType := textCmd.String("bbox-type", "xyxy", "Bbox type")
	alignY := textCmd.String("align-y", "baseline", "Y-axis alignment")
	alignX := textCmd.String("align-x", "left", "X-axis alignment")
	fontPath := textCmd.String("font", "", "Font file path")
	input := textCmd.String("input", "", "Input image path")
	imgz := textCmd.String("imgz", "1280x720", "Image size")
	output := textCmd.String("output", "", "Output image path")
	debug := textCmd.Bool("debug", false, "Enable debug output")

	help := textCmd.Bool("help", false, "Show help")
	version := textCmd.Bool("version", false, "Show version")

	textCmd.Parse(os.Args[2:])

	if *help {
		fmt.Print(Usage)
		return
	}

	if *version {
		fmt.Println(Version)
		return
	}

	if *output == "" {
		fmt.Println("Error: --output is required")
		os.Exit(1)
	}

	if *content == "" {
		fmt.Println("Error: --content is required")
		os.Exit(1)
	}

	// 解析字号
	sizeVal, sizeUnit, err := parseSizeValue(*size)
	if err != nil {
		fmt.Printf("Error parsing size: %v\n", err)
		os.Exit(1)
	}

	// 验证单位
	validUnits := map[string]bool{"px": true, "%": true, "rem": true, "bh": true}
	if !validUnits[sizeUnit] {
		fmt.Printf("Error: invalid size unit '%s'. Supported units: px, %%, rem, bh\n", sizeUnit)
		os.Exit(1)
	}

	// 解析颜色
	colorObj := parseColor(*colorStr)

	// 解析描边宽度
	strokeWidthVal := 0.0
	if *strokeWidth != "" {
		strokeWidthVal, _ = strconv.ParseFloat(*strokeWidth, 64)
	} else {
		strokeWidthVal, _ = strconv.ParseFloat(*stroke, 64)
	}

	// 保持向后兼容的逻辑
	strokeColorObj := colorObj
	if *strokeColorStr != "" {
		strokeColorObj = parseColor(*strokeColorStr)
		if *debug {
			fmt.Printf("Using custom stroke color: %s\n", *strokeColorStr)
		}
	} else if *debug && strokeWidthVal > 0 {
		fmt.Printf("Using text color for stroke (backward compatible)\n")
	}

	// 解析图片尺寸
	imgSize, err := parseImageSize(*imgz)
	if err != nil {
		fmt.Printf("Error parsing image size: %v\n", err)
		os.Exit(1)
	}

	// 解析边界框
	var bboxCoords [4]float64
	hasBBox := *bbox != ""
	if hasBBox {
		coords, err := parseBBox(*bbox, BBoxType(*bboxType))
		if err != nil {
			fmt.Printf("Error parsing bbox: %v\n", err)
			os.Exit(1)
		}
		bboxCoords = coords
	}

	// 解析对齐方式
	alignYObj, err := parseAlignment(*alignY, *rootSize)
	if err != nil {
		fmt.Printf("Error parsing align-y: %v\n", err)
		os.Exit(1)
	}

	alignXObj, err := parseAlignment(*alignX, *rootSize)
	if err != nil {
		fmt.Printf("Error parsing align-x: %v\n", err)
		os.Exit(1)
	}

	config := &TextConfig{
		Content:      *content,
		Size:         sizeVal,
		SizeUnit:     sizeUnit,
		RootSize:     *rootSize,
		Color:        colorObj,
		StrokeWidth:  strokeWidthVal,
		StrokeColor:  strokeColorObj,
		BBox:         bboxCoords,
		BBoxType:     BBoxType(*bboxType),
		AlignY:       alignYObj,
		AlignX:       alignXObj,
		FontPath:     *fontPath,
		InputPath:    *input,
		ImageSize:    imgSize,
		OutputPath:   *output,
		Debug:        *debug,
	}

	if err := runTextCommand(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runTextCommand(config *TextConfig) error {
	// 创建输出目录
	if err := os.MkdirAll(filepath.Dir(config.OutputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 加载或创建图像
	var img image.Image
	var err error

	if config.InputPath != "" {
		if config.Debug {
			fmt.Printf("Loading input image: %s\n", config.InputPath)
		}
		img, err = gg.LoadImage(config.InputPath)
		if err != nil {
			return fmt.Errorf("failed to load input image: %w", err)
		}
		if config.Debug {
			bounds := img.Bounds()
			fmt.Printf("Input image size: %dx%d\n", bounds.Dx(), bounds.Dy())
		}
	} else {
		if config.Debug {
			fmt.Printf("Creating transparent image: %dx%d\n", config.ImageSize.Width, config.ImageSize.Height)
		}
		img = image.NewRGBA(image.Rect(0, 0, config.ImageSize.Width, config.ImageSize.Height))
	}

	dc := gg.NewContextForImage(img)

	// 获取图像尺寸
	imgW := float64(dc.Width())
	imgH := float64(dc.Height())

	if config.Debug {
		fmt.Printf("Drawing context size: %.0fx%.0f\n", imgW, imgH)
	}

	// 确定边界框
	var bbox [4]float64
	if config.BBox != [4]float64{0, 0, 0, 0} {
		bbox = config.BBox
		if config.Debug {
			fmt.Printf("Using custom bbox: [%.0f, %.0f, %.0f, %.0f]\n", bbox[0], bbox[1], bbox[2], bbox[3])
		}
	} else {
		bbox = [4]float64{0, 0, imgW, imgH}
		if config.Debug {
			fmt.Printf("Using default bbox (full image): [%.0f, %.0f, %.0f, %.0f]\n", bbox[0], bbox[1], bbox[2], bbox[3])
		}
	}

	// 计算实际字体大小
	fontSize := config.Size
	bboxHeight := bbox[3] - bbox[1]

	switch config.SizeUnit {
	case "%":
		fontSize = bboxHeight * config.Size / 100.0
		if config.Debug {
			fmt.Printf("Font size from %%: %.0f%% of bbox height %.0f = %.0fpx\n",
				config.Size, bboxHeight, fontSize)
		}
	case "rem":
		fontSize = config.RootSize * config.Size
		if config.Debug {
			fmt.Printf("Font size from rem: %.2frem = %.0fpx (root-size: %.0fpx)\n",
				config.Size, fontSize, config.RootSize)
		}
	case "bh":
		fontSize = bboxHeight * config.Size
		if config.Debug {
			fmt.Printf("Font size from bh: %.2fbh = %.0fpx (bbox height: %.0fpx)\n",
				config.Size, fontSize, bboxHeight)
		}
	default:
		if config.Debug {
			fmt.Printf("Font size in px: %.0fpx\n", fontSize)
		}
	}

	// 加载字体
	fontFace, err := loadFont(config.FontPath, fontSize, config.Debug)
	if err != nil {
		return fmt.Errorf("failed to load font: %w", err)
	}
	dc.SetFontFace(fontFace)

	// 计算文本位置
	x, y := calculateTextPosition(dc, config, bbox)

	// 测量文本尺寸
	textW, textH := dc.MeasureString(config.Content)
	if config.Debug {
		fmt.Printf("Text: '%s'\n", config.Content)
		fmt.Printf("Text size: %.0f x %.0f\n", textW, textH)
		fmt.Printf("Text position: (%.0f, %.0f)\n", x, y)
		fmt.Printf("BBox bounds: X[%.0f-%.0f] Y[%.0f-%.0f]\n", bbox[0], bbox[2], bbox[1], bbox[3])

		if x < bbox[0] || x+textW > bbox[2] || y-textH < bbox[1] || y > bbox[3] {
			fmt.Printf("⚠️  Warning: Text may be outside bbox!\n")
		}
	}

	// 绘制描边
	if config.StrokeWidth > 0 {
		if config.Debug {
			if config.StrokeColor != config.Color {
				fmt.Printf("Drawing stroke with width: %.0f, color: %+v (custom)\n",
					config.StrokeWidth, config.StrokeColor)
			} else {
				fmt.Printf("Drawing stroke with width: %.0f, color: %+v (same as text)\n",
					config.StrokeWidth, config.StrokeColor)
			}
		}

		dc.Push()
		dc.SetColor(config.StrokeColor)

		strokeWidth := config.StrokeWidth
		offsets := []struct{ dx, dy float64 }{
			{-strokeWidth, 0}, {strokeWidth, 0},
			{0, -strokeWidth}, {0, strokeWidth},
		}

		if strokeWidth >= 1 {
			diagOffset := strokeWidth * 0.707
			offsets = append(offsets,
				struct{ dx, dy float64 }{-diagOffset, -diagOffset},
				struct{ dx, dy float64 }{diagOffset, diagOffset},
				struct{ dx, dy float64 }{-diagOffset, diagOffset},
				struct{ dx, dy float64 }{diagOffset, -diagOffset},
			)
		}

		for _, offset := range offsets {
			dc.DrawString(config.Content, x+offset.dx, y+offset.dy)
		}
		dc.Pop()
	}

	// 绘制主文本
	if config.Debug {
		fmt.Printf("Drawing main text with color: %+v at (%.0f, %.0f)\n", config.Color, x, y)
	}
	dc.SetColor(config.Color)
	dc.DrawString(config.Content, x, y)

	// 保存图像
	if err := dc.SavePNG(config.OutputPath); err != nil {
		return fmt.Errorf("failed to save output image: %w", err)
	}

	fmt.Printf("✅ Watermark text added successfully: %s\n", config.OutputPath)
	return nil
}

// loadFont 加载字体
func loadFont(fontPath string, fontSize float64, debug bool) (font.Face, error) {
	if fontPath == "" {
		if debug {
			fmt.Println("Using built-in font (English only)")
		}
		return loadBuiltinFont(fontSize)
	}

	ext := strings.ToLower(filepath.Ext(fontPath))
	if ext == ".ttf" || ext == ".ttc" || ext == ".otf" {
		if debug {
			fmt.Printf("Loading font from file: %s\n", fontPath)
		}
		return loadFontFromFile(fontPath, fontSize)
	}

	if debug {
		fmt.Printf("Searching system font by name: %s\n", fontPath)
	}
	return loadFontByName(fontPath, fontSize, debug)
}

// loadBuiltinFont 加载内置字体
func loadBuiltinFont(fontSize float64) (font.Face, error) {
	parsedFont, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse builtin font: %w", err)
	}

	fontFace, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}

	return fontFace, nil
}

// loadFontFromFile 从文件加载字体（支持 TTF 和 TTC）
func loadFontFromFile(path string, fontSize float64) (font.Face, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("font file not found: %w", err)
	}

	fontData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %w", err)
	}

	var parsedFont *opentype.Font

	// 尝试直接解析
	parsedFont, err = opentype.Parse(fontData)
	if err != nil {
		// 如果是 TTC，尝试解析集合
		if strings.HasSuffix(strings.ToLower(path), ".ttc") {
			collection, collectionErr := opentype.ParseCollection(fontData)
			if collectionErr != nil {
				return nil, fmt.Errorf("failed to parse TTC collection: %w", collectionErr)
			}
			if collection.NumFonts() == 0 {
				return nil, fmt.Errorf("TTC collection contains no fonts")
			}
			parsedFont, err = collection.Font(0)
			if err != nil {
				return nil, fmt.Errorf("failed to get font from collection: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse font: %w", err)
		}
	}

	fontFace, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}

	return fontFace, nil
}

// loadFontByName 按字体名从系统查找并加载
func loadFontByName(fontName string, fontSize float64, debug bool) (font.Face, error) {
	cmd := exec.Command("fc-list", "-f", "%{file}", fontName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to search font '%s': %w", fontName, err)
	}

	fontPath := strings.TrimSpace(string(output))
	if fontPath == "" {
		return nil, fmt.Errorf("font '%s' not found in system", fontName)
	}

	if debug {
		fmt.Printf("Found font file: %s\n", fontPath)
	}

	return loadFontFromFile(fontPath, fontSize)
}

// calculateTextPosition 计算文本位置
func calculateTextPosition(dc *gg.Context, config *TextConfig, bbox [4]float64) (float64, float64) {
	textW, textH := dc.MeasureString(config.Content)

	x := bbox[0]
	switch config.AlignX.Type {
	case "right":
		x = bbox[2] - textW
	case "middle", "center":
		x = bbox[0] + (bbox[2]-bbox[0])/2 - textW/2
	case "left":
		x = bbox[0]
	default:
		x = bbox[0]
	}

	if config.AlignX.Offset != 0 {
		x += config.AlignX.Offset
	}
	if config.AlignX.IsRel && config.AlignX.OffsetRel != 0 {
		x += (bbox[2] - bbox[0]) * config.AlignX.OffsetRel
	}

	y := bbox[1]
	switch config.AlignY.Type {
	case "top":
		y = bbox[1] + textH
	case "bottom":
		y = bbox[3]
	case "middle", "center":
		y = bbox[1] + (bbox[3]-bbox[1])/2 + textH/2
	case "baseline":
		y = bbox[1] + textH
	default:
		y = bbox[1] + textH
	}

	if config.AlignY.Offset != 0 {
		y += config.AlignY.Offset
	}
	if config.AlignY.IsRel && config.AlignY.OffsetRel != 0 {
		y += (bbox[3] - bbox[1]) * config.AlignY.OffsetRel
	}

	return x, y
}

// parseSizeValue 解析大小值
func parseSizeValue(sizeStr string) (float64, string, error) {
	sizeStr = strings.TrimSpace(sizeStr)

	units := []string{"px", "%", "rem", "bh"}
	for _, unit := range units {
		if strings.HasSuffix(sizeStr, unit) {
			valStr := strings.TrimSuffix(sizeStr, unit)
			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return 0, "", err
			}
			return val, unit, nil
		}
	}

	val, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, "", err
	}
	return val, "px", nil
}

// parseColor 解析颜色
func parseColor(colorStr string) color.Color {
	colorStr = strings.ToLower(strings.TrimSpace(colorStr))

	switch colorStr {
	case "black":
		return color.RGBA{0, 0, 0, 255}
	case "white":
		return color.RGBA{255, 255, 255, 255}
	case "red":
		return color.RGBA{255, 0, 0, 255}
	case "green":
		return color.RGBA{0, 255, 0, 255}
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "yellow":
		return color.RGBA{255, 255, 0, 255}
	case "cyan":
		return color.RGBA{0, 255, 255, 255}
	case "magenta":
		return color.RGBA{255, 0, 255, 255}
	case "transparent":
		return color.RGBA{0, 0, 0, 0}
	}

	if strings.HasPrefix(colorStr, "#") {
		var r, g, b, a uint8 = 0, 0, 0, 255
		switch len(colorStr) {
		case 7:
			fmt.Sscanf(colorStr, "#%02x%02x%02x", &r, &g, &b)
		case 9:
			fmt.Sscanf(colorStr, "#%02x%02x%02x%02x", &r, &g, &b, &a)
		}
		return color.RGBA{r, g, b, a}
	}

	if strings.HasPrefix(colorStr, "rgba(") && strings.HasSuffix(colorStr, ")") {
		inner := strings.TrimPrefix(colorStr, "rgba(")
		inner = strings.TrimSuffix(inner, ")")
		parts := strings.Split(inner, ",")
		if len(parts) == 4 {
			r, _ := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 8)
			g, _ := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 8)
			b, _ := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 8)
			a, _ := strconv.ParseUint(strings.TrimSpace(parts[3]), 10, 8)
			return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		}
	}

	return color.RGBA{255, 255, 255, 255}
}

// parseImageSize 解析图片尺寸
func parseImageSize(sizeStr string) (Size, error) {
	sizeStr = strings.TrimSpace(sizeStr)

	dim := strings.Split(sizeStr, "x")
	if len(dim) != 2 {
		return Size{}, fmt.Errorf("invalid size format, expected WIDTHxHEIGHT")
	}

	width, err := strconv.Atoi(dim[0])
	if err != nil {
		return Size{}, err
	}

	height, err := strconv.Atoi(dim[1])
	if err != nil {
		return Size{}, err
	}

	if width <= 0 || height <= 0 {
		return Size{}, fmt.Errorf("width and height must be positive")
	}

	return Size{Width: width, Height: height}, nil
}

// parseBBox 解析边界框
func parseBBox(bboxStr string, bboxType BBoxType) ([4]float64, error) {
	parts := strings.Split(bboxStr, ",")
	if len(parts) != 4 {
		return [4]float64{}, fmt.Errorf("bbox requires 4 values")
	}

	coords := [4]float64{}
	for i, part := range parts {
		val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return [4]float64{}, err
		}
		coords[i] = val
	}

	switch bboxType {
	case BBoxXYWH:
		return [4]float64{
			coords[0],
			coords[1],
			coords[0] + coords[2],
			coords[1] + coords[3],
		}, nil
	case BBoxCXCY:
		cx, cy, w, h := coords[0], coords[1], coords[2], coords[3]
		return [4]float64{
			cx - w/2,
			cy - h/2,
			cx + w/2,
			cy + h/2,
		}, nil
	default:
		return coords, nil
	}
}

// parseAlignment 解析对齐方式
func parseAlignment(alignStr string, rootSize float64) (Alignment, error) {
	alignStr = strings.TrimSpace(alignStr)

	offset := 0.0
	isRel := false
	relValue := 0.0

	plusIdx := strings.Index(alignStr, "+")
	minusIdx := strings.Index(alignStr, "-")

	var alignType string
	var offsetStr string

	if plusIdx > 0 {
		alignType = alignStr[:plusIdx]
		offsetStr = alignStr[plusIdx+1:]
	} else if minusIdx > 0 {
		alignType = alignStr[:minusIdx]
		offsetStr = alignStr[minusIdx:]
	} else {
		alignType = alignStr
	}

	if offsetStr != "" {
		if strings.HasSuffix(offsetStr, "%") {
			isRel = true
			val, err := strconv.ParseFloat(strings.TrimSuffix(offsetStr, "%"), 64)
			if err != nil {
				return Alignment{}, err
			}
			relValue = val / 100.0
		} else if strings.HasSuffix(offsetStr, "rem") {
			val, err := strconv.ParseFloat(strings.TrimSuffix(offsetStr, "rem"), 64)
			if err != nil {
				return Alignment{}, err
			}
			offset = val * rootSize
		} else if strings.HasSuffix(offsetStr, "px") {
			val, err := strconv.ParseFloat(strings.TrimSuffix(offsetStr, "px"), 64)
			if err != nil {
				return Alignment{}, err
			}
			offset = val
		} else {
			val, err := strconv.ParseFloat(offsetStr, 64)
			if err != nil {
				return Alignment{}, err
			}
			offset = val
		}
	}

	return Alignment{
		Type:      alignType,
		Offset:    offset,
		OffsetRel: relValue,
		IsRel:     isRel,
	}, nil
}