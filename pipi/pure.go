// Package pipi provides pure functions for text watermark processing
package pipi

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// ============================================================================
// Type Definitions
// ============================================================================

// BBoxType represents bounding box coordinate type
type BBoxType string

const (
	BBoxXYXY BBoxType = "xyxy" // x1,y1,x2,y2
	BBoxXYWH BBoxType = "xywh" // x,y,width,height
	BBoxCXCY BBoxType = "cxcy" // center-x,center-y,width,height
)

// FontSizeUnit represents font size unit types
type FontSizeUnit string

const (
	UnitPX  FontSizeUnit = "px"  // absolute pixels
	UnitPct FontSizeUnit = "%"   // percentage of bbox height
	UnitREM FontSizeUnit = "rem" // relative to root size
	UnitBH  FontSizeUnit = "bh"  // box height multiplier
)

// Alignment represents text alignment configuration
type Alignment struct {
	Type      string  // baseline, top, bottom, middle, left, right, center
	Offset    float64 // pixel offset
	OffsetRel float64 // relative offset (0-1)
	IsRel     bool    // whether using relative offset
}

// StrokeStyle represents text stroke configuration
type StrokeStyle struct {
	Width float64
	Color color.Color
}

// TextMetrics contains measured text dimensions
type TextMetrics struct {
	Width  float64
	Height float64
}

// BBox represents a bounding box with four coordinates
type BBox [4]float64

// ============================================================================
// Size and Unit Parsing
// ============================================================================

// ParseSizeValue parses a size string like "24px", "2rem", "10%", "0.5bh"
// Returns the numeric value, unit, and any error encountered.
func ParseSizeValue(sizeStr string) (float64, string, error) {
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

	// Default to px if no unit specified
	val, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, "", err
	}
	return val, "px", nil
}

// ParseImageSize parses image size string like "1280x720"
// Returns width, height, and any error encountered.
func ParseImageSize(sizeStr string) (width, height int, err error) {
	sizeStr = strings.TrimSpace(sizeStr)

	dims := strings.Split(sizeStr, "x")
	if len(dims) != 2 {
		return 0, 0, fmt.Errorf("invalid size format, expected WIDTHxHEIGHT")
	}

	width, err = strconv.Atoi(strings.TrimSpace(dims[0]))
	if err != nil {
		return 0, 0, err
	}

	height, err = strconv.Atoi(strings.TrimSpace(dims[1]))
	if err != nil {
		return 0, 0, err
	}

	if width <= 0 || height <= 0 {
		return 0, 0, fmt.Errorf("width and height must be positive")
	}

	return width, height, nil
}

// IsValidSizeUnit checks if the given unit is valid
func IsValidSizeUnit(unit string) bool {
	validUnits := map[string]bool{"px": true, "%": true, "rem": true, "bh": true}
	return validUnits[unit]
}

// ============================================================================
// Color Processing
// ============================================================================

// ParseColor parses a color string, supporting color names, hex, and rgba()
// Returns a color.Color or default white if parsing fails.
func ParseColor(colorStr string) color.Color {
	return ParseColorWithDefault(colorStr, color.RGBA{255, 255, 255, 255})
}

// ParseColorWithDefault parses a color string, returning defaultColor on failure
func ParseColorWithDefault(colorStr string, defaultColor color.Color) color.Color {
	colorStr = strings.ToLower(strings.TrimSpace(colorStr))

	// Named colors
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

	// Hex color
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

	// RGBA color
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

	return defaultColor
}

// ColorToRGBA converts color.Color to RGBA components
func ColorToRGBA(c color.Color) (r, g, b, a uint8) {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return rgba.R, rgba.G, rgba.B, rgba.A
}

// ============================================================================
// Bounding Box Processing
// ============================================================================

// ParseBBox parses a bounding box string based on the specified type
// Returns x1, y1, x2, y2 coordinates and any error encountered.
func ParseBBox(bboxStr string, bboxType BBoxType) (x1, y1, x2, y2 float64, err error) {
	parts := strings.Split(bboxStr, ",")
	if len(parts) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("bbox requires 4 values")
	}

	coords := [4]float64{}
	for i, part := range parts {
		val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return 0, 0, 0, 0, err
		}
		coords[i] = val
	}

	switch bboxType {
	case BBoxXYWH:
		return coords[0], coords[1], coords[0] + coords[2], coords[1] + coords[3], nil
	case BBoxCXCY:
		cx, cy, w, h := coords[0], coords[1], coords[2], coords[3]
		return cx - w/2, cy - h/2, cx + w/2, cy + h/2, nil
	default: // BBoxXYXY
		return coords[0], coords[1], coords[2], coords[3], nil
	}
}

// ConvertBBox converts bounding box between different formats
func ConvertBBox(x1, y1, x2, y2 float64, toType BBoxType) [4]float64 {
	width := x2 - x1
	height := y2 - y1

	switch toType {
	case BBoxXYWH:
		return [4]float64{x1, y1, width, height}
	case BBoxCXCY:
		return [4]float64{x1 + width/2, y1 + height/2, width, height}
	default: // BBoxXYXY
		return [4]float64{x1, y1, x2, y2}
	}
}

// ValidateBBox validates if a bounding box is valid (x2 > x1, y2 > y1)
func ValidateBBox(x1, y1, x2, y2 float64) bool {
	return x2 > x1 && y2 > y1
}

// GetDefaultBBox returns the default bounding box for an image
func GetDefaultBBox(imgWidth, imgHeight float64) BBox {
	return BBox{0, 0, imgWidth, imgHeight}
}

// ============================================================================
// Font Size Calculation
// ============================================================================

// CalculateFontSize calculates actual font size based on unit and context
func CalculateFontSize(sizeValue float64, unit FontSizeUnit, rootSize, bboxHeight float64) float64 {
	switch unit {
	case UnitPct:
		return bboxHeight * sizeValue / 100.0
	case UnitREM:
		return rootSize * sizeValue
	case UnitBH:
		return bboxHeight * sizeValue
	default: // UnitPX
		return sizeValue
	}
}

// ParseAndCalculateFontSize parses size string and calculates actual font size
func ParseAndCalculateFontSize(sizeStr string, rootSize, bboxHeight float64) (fontSize float64, unit FontSizeUnit, err error) {
	value, unitStr, err := ParseSizeValue(sizeStr)
	if err != nil {
		return 0, "", err
	}

	if !IsValidSizeUnit(unitStr) {
		return 0, "", fmt.Errorf("invalid size unit: %s", unitStr)
	}

	fontSize = CalculateFontSize(value, FontSizeUnit(unitStr), rootSize, bboxHeight)
	return fontSize, FontSizeUnit(unitStr), nil
}

// ============================================================================
// Alignment Processing
// ============================================================================

// ParseAlignment parses alignment string like "bottom+10px", "middle-5%"
// Supports pixel, rem, and percentage offsets.
func ParseAlignment(alignStr string, rootSize float64) (Alignment, error) {
	alignStr = strings.TrimSpace(alignStr)

	offset := 0.0
	isRel := false
	relValue := 0.0

	// Find offset separator
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

	// Parse offset if present
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

// CalculateAlignedPositionX calculates X position based on horizontal alignment
func CalculateAlignedPositionX(bboxMinX, bboxMaxX, textWidth float64, align Alignment) float64 {
	var x float64

	switch align.Type {
	case "right":
		x = bboxMaxX - textWidth
	case "middle", "center":
		x = bboxMinX + (bboxMaxX-bboxMinX)/2 - textWidth/2
	case "left":
		x = bboxMinX
	default:
		x = bboxMinX
	}

	// Apply offsets
	if align.Offset != 0 {
		x += align.Offset
	}
	if align.IsRel && align.OffsetRel != 0 {
		x += (bboxMaxX - bboxMinX) * align.OffsetRel
	}

	return x
}

// CalculateAlignedPositionY calculates Y position based on vertical alignment
func CalculateAlignedPositionY(bboxMinY, bboxMaxY, textHeight float64, align Alignment) float64 {
	y := bboxMinY

	switch align.Type {
	case "top":
		y = bboxMinY + textHeight
	case "bottom":
		y = bboxMaxY
	case "middle", "center":
		y = bboxMinY + (bboxMaxY-bboxMinY)/2 + textHeight/2
	case "baseline":
		y = bboxMinY + textHeight
	default:
		y = bboxMinY + textHeight
	}

	// Apply offsets
	if align.Offset != 0 {
		y += align.Offset
	}
	if align.IsRel && align.OffsetRel != 0 {
		y += (bboxMaxY - bboxMinY) * align.OffsetRel
	}

	return y
}

// CalculateAlignedPosition calculates both X and Y positions
func CalculateAlignedPosition(bbox BBox, textWidth, textHeight float64, alignX, alignY Alignment) (x, y float64) {
	x = CalculateAlignedPositionX(bbox[0], bbox[2], textWidth, alignX)
	y = CalculateAlignedPositionY(bbox[1], bbox[3], textHeight, alignY)
	return x, y
}

// ============================================================================
// Text Bounds Checking
// ============================================================================

// CheckTextBounds checks if text rectangle is within bounding box
func CheckTextBounds(x, y, textWidth, textHeight float64, bbox BBox) bool {
	// Text is drawn from (x, y) where y is the baseline
	// The text extends upward from baseline by approx textHeight
	textTop := y - textHeight
	textBottom := y
	textLeft := x
	textRight := x + textWidth

	return textLeft >= bbox[0] && textRight <= bbox[2] &&
		textTop >= bbox[1] && textBottom <= bbox[3]
}

// GetTextBounds returns the bounding box of rendered text
func GetTextBounds(x, y, textWidth, textHeight float64) BBox {
	return BBox{
		x,              // left
		y - textHeight, // top
		x + textWidth,  // right
		y,              // bottom (baseline)
	}
}

// ============================================================================
// Stroke Processing
// ============================================================================

// GenerateStrokeOffsets generates offset vectors for stroke rendering
// Returns a slice of (dx, dy) offset pairs
func GenerateStrokeOffsets(strokeWidth float64) []struct{ Dx, Dy float64 } {
	offsets := []struct{ Dx, Dy float64 }{
		{-strokeWidth, 0}, {strokeWidth, 0},
		{0, -strokeWidth}, {0, strokeWidth},
	}

	// Add diagonal offsets for thicker strokes
	if strokeWidth >= 1 {
		diagOffset := strokeWidth * 0.707 // sqrt(2)/2
		offsets = append(offsets,
			struct{ Dx, Dy float64 }{-diagOffset, -diagOffset},
			struct{ Dx, Dy float64 }{diagOffset, diagOffset},
			struct{ Dx, Dy float64 }{-diagOffset, diagOffset},
			struct{ Dx, Dy float64 }{diagOffset, -diagOffset},
		)
	}

	return offsets
}

// GetStrokeColor determines the stroke color to use
func GetStrokeColor(strokeWidth float64, strokeColor, textColor color.Color) color.Color {
	if strokeWidth > 0 && strokeColor != nil {
		return strokeColor
	}
	return textColor
}

// ============================================================================
// Configuration Validation
// ============================================================================

// ValidateTextConfig validates text configuration parameters
func ValidateTextConfig(content string, fontSize float64, bbox BBox) error {
	if content == "" {
		return fmt.Errorf("text content cannot be empty")
	}

	if fontSize <= 0 {
		return fmt.Errorf("font size must be positive, got %f", fontSize)
	}

	if !ValidateBBox(bbox[0], bbox[1], bbox[2], bbox[3]) {
		return fmt.Errorf("invalid bounding box: [%v, %v, %v, %v]", bbox[0], bbox[1], bbox[2], bbox[3])
	}

	return nil
}

// ============================================================================
// String Utilities
// ============================================================================

// NormalizeCommand normalizes command string by handling aliases
func NormalizeCommand(cmd string) string {
	switch cmd {
	case "center":
		return "middle"
	default:
		return cmd
	}
}

// ParseCoordinate parses a coordinate value that may have unit suffixes
func ParseCoordinate(coordStr string, rootSize float64) (float64, error) {
	coordStr = strings.TrimSpace(coordStr)

	if strings.HasSuffix(coordStr, "px") {
		return strconv.ParseFloat(strings.TrimSuffix(coordStr, "px"), 64)
	}

	if strings.HasSuffix(coordStr, "rem") {
		val, err := strconv.ParseFloat(strings.TrimSuffix(coordStr, "rem"), 64)
		if err != nil {
			return 0, err
		}
		return val * rootSize, nil
	}

	if strings.HasSuffix(coordStr, "%") {
		return strconv.ParseFloat(strings.TrimSuffix(coordStr, "%"), 64)
	}

	return strconv.ParseFloat(coordStr, 64)
}