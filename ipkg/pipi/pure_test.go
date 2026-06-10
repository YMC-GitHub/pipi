package pipi

import (
	"image/color"
	"testing"
)

// ============================================================================
// Size and Unit Parsing Tests
// ============================================================================

func TestParseSizeValue(t *testing.T) {
	tests := []struct {
		name      string
		sizeStr   string
		wantValue float64
		wantUnit  string
		wantErr   bool
	}{
		{"px unit", "24px", 24, "px", false},
		{"rem unit", "2rem", 2, "rem", false},
		{"percentage unit", "10%", 10, "%", false},
		{"bh unit", "0.5bh", 0.5, "bh", false},
		{"no unit (default px)", "24", 24, "px", false},
		{"decimal value", "1.5rem", 1.5, "rem", false},
		{"negative value", "-5px", -5, "px", false},
		{"invalid number", "abcpx", 0, "", true},
		{"empty string", "", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotUnit, err := ParseSizeValue(tt.sizeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSizeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotValue != tt.wantValue {
					t.Errorf("ParseSizeValue() value = %v, want %v", gotValue, tt.wantValue)
				}
				if gotUnit != tt.wantUnit {
					t.Errorf("ParseSizeValue() unit = %v, want %v", gotUnit, tt.wantUnit)
				}
			}
		})
	}
}

func TestParseImageSize(t *testing.T) {
	tests := []struct {
		name       string
		sizeStr    string
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{"valid size", "1280x720", 1280, 720, false},
		{"small size", "100x100", 100, 100, false},
		{"with spaces", " 800x600 ", 800, 600, false},
		{"invalid format", "1280", 0, 0, true},
		{"invalid format no x", "1280*720", 0, 0, true},
		{"negative width", "-100x200", 0, 0, true},
		{"negative height", "100x-200", 0, 0, true},
		{"zero width", "0x100", 0, 0, true},
		{"non-numeric", "abxc", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight, err := ParseImageSize(tt.sizeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseImageSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotWidth != tt.wantWidth {
					t.Errorf("ParseImageSize() width = %v, want %v", gotWidth, tt.wantWidth)
				}
				if gotHeight != tt.wantHeight {
					t.Errorf("ParseImageSize() height = %v, want %v", gotHeight, tt.wantHeight)
				}
			}
		})
	}
}

func TestIsValidSizeUnit(t *testing.T) {
	tests := []struct {
		unit string
		want bool
	}{
		{"px", true},
		{"%", true},
		{"rem", true},
		{"bh", true},
		{"em", false},
		{"pt", false},
		{"", false},
		{"Px", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			if got := IsValidSizeUnit(tt.unit); got != tt.want {
				t.Errorf("IsValidSizeUnit(%q) = %v, want %v", tt.unit, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Color Processing Tests
// ============================================================================

func TestParseColor(t *testing.T) {
	tests := []struct {
		name      string
		colorStr  string
		wantColor color.Color
	}{
		{"named color black", "black", color.RGBA{0, 0, 0, 255}},
		{"named color white", "white", color.RGBA{255, 255, 255, 255}},
		{"named color red", "red", color.RGBA{255, 0, 0, 255}},
		{"named color green", "green", color.RGBA{0, 255, 0, 255}},
		{"named color blue", "blue", color.RGBA{0, 0, 255, 255}},
		{"named color yellow", "yellow", color.RGBA{255, 255, 0, 255}},
		{"named color cyan", "cyan", color.RGBA{0, 255, 255, 255}},
		{"named color magenta", "magenta", color.RGBA{255, 0, 255, 255}},
		{"named color transparent", "transparent", color.RGBA{0, 0, 0, 0}},
		{"hex 6-digit", "#FF0000", color.RGBA{255, 0, 0, 255}},
		{"hex 6-digit lowercase", "#ff0000", color.RGBA{255, 0, 0, 255}},
		{"hex 8-digit", "#FF0000FF", color.RGBA{255, 0, 0, 255}},
		{"rgba", "rgba(255,0,0,128)", color.RGBA{255, 0, 0, 128}},
		{"rgba with spaces", "rgba( 255 , 0 , 0 , 128 )", color.RGBA{255, 0, 0, 128}},
		{"unknown color defaults to white", "unknown", color.RGBA{255, 255, 255, 255}},
		{"empty string defaults to white", "", color.RGBA{255, 255, 255, 255}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseColor(tt.colorStr)
			if !colorsEqual(got, tt.wantColor) {
				t.Errorf("ParseColor(%q) = %v, want %v", tt.colorStr, got, tt.wantColor)
			}
		})
	}
}

func TestParseColorWithDefault(t *testing.T) {
	defaultColor := color.RGBA{128, 128, 128, 255}

	tests := []struct {
		name     string
		colorStr string
		want     color.Color
	}{
		{"valid color", "red", color.RGBA{255, 0, 0, 255}},
		{"invalid color uses default", "invalid", defaultColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseColorWithDefault(tt.colorStr, defaultColor)
			if !colorsEqual(got, tt.want) {
				t.Errorf("ParseColorWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorToRGBA(t *testing.T) {
	tests := []struct {
		name  string
		color color.Color
	}{
		{"RGBA color", color.RGBA{255, 128, 64, 200}},
		{"NRGBA color", color.NRGBA{100, 150, 200, 180}},
		{"Gray color", color.Gray{100}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, a := ColorToRGBA(tt.color)
			// Convert to RGBA for consistent comparison
			rgba := color.RGBAModel.Convert(tt.color).(color.RGBA)
			if r != rgba.R || g != rgba.G || b != rgba.B || a != rgba.A {
				t.Errorf("ColorToRGBA() = (%d,%d,%d,%d), want (%d,%d,%d,%d)",
					r, g, b, a, rgba.R, rgba.G, rgba.B, rgba.A)
			}
		})
	}
}

// ============================================================================
// Bounding Box Tests
// ============================================================================

func TestParseBBox(t *testing.T) {
	tests := []struct {
		name     string
		bboxStr  string
		bboxType BBoxType
		wantX1   float64
		wantY1   float64
		wantX2   float64
		wantY2   float64
		wantErr  bool
	}{
		{"xyxy format", "10,20,100,200", BBoxXYXY, 10, 20, 100, 200, false},
		{"xywh format", "10,20,90,180", BBoxXYWH, 10, 20, 100, 200, false},
		{"cxcy format", "55,110,90,180", BBoxCXCY, 10, 20, 100, 200, false},
		{"with spaces", " 10 , 20 , 100 , 200 ", BBoxXYXY, 10, 20, 100, 200, false},
		{"invalid number of values", "10,20,100", BBoxXYXY, 0, 0, 0, 0, true},
		{"invalid number", "10,abc,100,200", BBoxXYXY, 0, 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX1, gotY1, gotX2, gotY2, err := ParseBBox(tt.bboxStr, tt.bboxType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBBox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotX1 != tt.wantX1 || gotY1 != tt.wantY1 || gotX2 != tt.wantX2 || gotY2 != tt.wantY2 {
					t.Errorf("ParseBBox() = (%v,%v,%v,%v), want (%v,%v,%v,%v)",
						gotX1, gotY1, gotX2, gotY2, tt.wantX1, tt.wantY1, tt.wantX2, tt.wantY2)
				}
			}
		})
	}
}

func TestConvertBBox(t *testing.T) {
	tests := []struct {
		name   string
		x1, y1 float64
		x2, y2 float64
		toType BBoxType
		want   [4]float64
	}{
		{"XYXY to XYWH", 10, 20, 100, 200, BBoxXYWH, [4]float64{10, 20, 90, 180}},
		{"XYXY to CXCY", 10, 20, 100, 200, BBoxCXCY, [4]float64{55, 110, 90, 180}},
		{"XYXY to XYXY", 10, 20, 100, 200, BBoxXYXY, [4]float64{10, 20, 100, 200}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertBBox(tt.x1, tt.y1, tt.x2, tt.y2, tt.toType)
			if got != tt.want {
				t.Errorf("ConvertBBox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateBBox(t *testing.T) {
	tests := []struct {
		name   string
		x1, y1 float64
		x2, y2 float64
		want   bool
	}{
		{"valid", 0, 0, 100, 100, true},
		{"invalid x", 100, 0, 0, 100, false},
		{"invalid y", 0, 100, 100, 0, false},
		{"zero width", 0, 0, 0, 100, false},
		{"zero height", 0, 0, 100, 0, false},
		{"negative values", -10, -10, 10, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateBBox(tt.x1, tt.y1, tt.x2, tt.y2); got != tt.want {
				t.Errorf("ValidateBBox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultBBox(t *testing.T) {
	tests := []struct {
		width  float64
		height float64
		want   BBox
	}{
		{100, 200, BBox{0, 0, 100, 200}},
		{0, 100, BBox{0, 0, 0, 100}},
		{100, 0, BBox{0, 0, 100, 0}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := GetDefaultBBox(tt.width, tt.height)
			if got != tt.want {
				t.Errorf("GetDefaultBBox() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Font Size Calculation Tests
// ============================================================================

func TestCalculateFontSize(t *testing.T) {
	tests := []struct {
		name       string
		sizeValue  float64
		unit       FontSizeUnit
		rootSize   float64
		bboxHeight float64
		want       float64
	}{
		{"px unit", 24, UnitPX, 16, 100, 24},
		{"rem unit", 2, UnitREM, 16, 100, 32},
		{"percentage", 50, UnitPct, 16, 100, 50},
		{"bh unit", 0.5, UnitBH, 16, 100, 50},
		{"zero value", 0, UnitPX, 16, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateFontSize(tt.sizeValue, tt.unit, tt.rootSize, tt.bboxHeight)
			if got != tt.want {
				t.Errorf("CalculateFontSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAndCalculateFontSize(t *testing.T) {
	tests := []struct {
		name       string
		sizeStr    string
		rootSize   float64
		bboxHeight float64
		wantSize   float64
		wantUnit   FontSizeUnit
		wantErr    bool
	}{
		{"px string", "24px", 16, 100, 24, UnitPX, false},
		{"rem string", "2rem", 16, 100, 32, UnitREM, false},
		{"percentage", "50%", 16, 100, 50, UnitPct, false},
		{"bh string", "0.5bh", 16, 100, 50, UnitBH, false},
		{"invalid unit", "24em", 16, 100, 0, "", true},
		{"invalid number", "abcp", 16, 100, 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSize, gotUnit, err := ParseAndCalculateFontSize(tt.sizeStr, tt.rootSize, tt.bboxHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAndCalculateFontSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotSize != tt.wantSize {
					t.Errorf("ParseAndCalculateFontSize() size = %v, want %v", gotSize, tt.wantSize)
				}
				if gotUnit != tt.wantUnit {
					t.Errorf("ParseAndCalculateFontSize() unit = %v, want %v", gotUnit, tt.wantUnit)
				}
			}
		})
	}
}

// ============================================================================
// Alignment Tests
// ============================================================================

func TestParseAlignment(t *testing.T) {
	tests := []struct {
		name     string
		alignStr string
		rootSize float64
		want     Alignment
		wantErr  bool
	}{
		{"simple baseline", "baseline", 16, Alignment{Type: "baseline", Offset: 0, OffsetRel: 0, IsRel: false}, false},
		{"bottom with px offset", "bottom+10px", 16, Alignment{Type: "bottom", Offset: 10, OffsetRel: 0, IsRel: false}, false},
		{"middle with negative px", "middle-5px", 16, Alignment{Type: "middle", Offset: -5, OffsetRel: 0, IsRel: false}, false},
		{"top with rem offset", "top+2rem", 16, Alignment{Type: "top", Offset: 32, OffsetRel: 0, IsRel: false}, false},
		{"center with percent", "center-10%", 16, Alignment{Type: "center", Offset: 0, OffsetRel: -0.1, IsRel: true}, false},
		{"left with percent", "left+25%", 16, Alignment{Type: "left", Offset: 0, OffsetRel: 0.25, IsRel: true}, false},
		{"right with rem", "right-1.5rem", 16, Alignment{Type: "right", Offset: -24, OffsetRel: 0, IsRel: false}, false},
		{"no type", "", 16, Alignment{Type: "", Offset: 0, OffsetRel: 0, IsRel: false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAlignment(tt.alignStr, tt.rootSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAlignment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Type != tt.want.Type {
					t.Errorf("ParseAlignment() type = %v, want %v", got.Type, tt.want.Type)
				}
				if got.Offset != tt.want.Offset {
					t.Errorf("ParseAlignment() offset = %v, want %v", got.Offset, tt.want.Offset)
				}
				if got.OffsetRel != tt.want.OffsetRel {
					t.Errorf("ParseAlignment() offsetRel = %v, want %v", got.OffsetRel, tt.want.OffsetRel)
				}
				if got.IsRel != tt.want.IsRel {
					t.Errorf("ParseAlignment() isRel = %v, want %v", got.IsRel, tt.want.IsRel)
				}
			}
		})
	}
}

func TestCalculateAlignedPositionX(t *testing.T) {
	bboxMinX, bboxMaxX := 100.0, 500.0
	textWidth := 200.0

	tests := []struct {
		name  string
		align Alignment
		want  float64
	}{
		{"left align", Alignment{Type: "left"}, 100},
		{"right align", Alignment{Type: "right"}, 300},
		{"center align", Alignment{Type: "center"}, 200},  // Fixed: 200 is correct
		{"middle align", Alignment{Type: "middle"}, 200}, // Fixed: 200 is correct
		{"default", Alignment{Type: "unknown"}, 100},
		{"left with offset", Alignment{Type: "left", Offset: 10}, 110},
		{"right with offset", Alignment{Type: "right", Offset: -10}, 290},
		{"center with relative offset", Alignment{Type: "center", IsRel: true, OffsetRel: 0.1}, 240}, // Fixed: 200 + 400*0.1 = 240
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAlignedPositionX(bboxMinX, bboxMaxX, textWidth, tt.align)
			if got != tt.want {
				t.Errorf("CalculateAlignedPositionX() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestCalculateAlignedPositionY(t *testing.T) {
	bboxMinY, bboxMaxY := 100.0, 500.0
	textHeight := 50.0

	tests := []struct {
		name  string
		align Alignment
		want  float64
	}{
		{"baseline", Alignment{Type: "baseline"}, 150},
		{"top", Alignment{Type: "top"}, 150},
		{"bottom", Alignment{Type: "bottom"}, 500},
		{"center", Alignment{Type: "center"}, 325},
		{"middle", Alignment{Type: "middle"}, 325},
		{"default", Alignment{Type: "unknown"}, 150},
		{"top with offset", Alignment{Type: "top", Offset: 10}, 160},
		{"bottom with offset", Alignment{Type: "bottom", Offset: -10}, 490},
		{"center with relative offset", Alignment{Type: "center", IsRel: true, OffsetRel: 0.1}, 365},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAlignedPositionY(bboxMinY, bboxMaxY, textHeight, tt.align)
			if got != tt.want {
				t.Errorf("CalculateAlignedPositionY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateAlignedPosition(t *testing.T) {
	bbox := BBox{100, 100, 500, 500}
	textWidth, textHeight := 200.0, 50.0
	alignX := Alignment{Type: "center"}
	alignY := Alignment{Type: "middle"}

	x, y := CalculateAlignedPosition(bbox, textWidth, textHeight, alignX, alignY)

	expectedX := 200.0  // Fixed: 200 is correct
	expectedY := 325.0

	if x != expectedX {
		t.Errorf("CalculateAlignedPosition() x = %v, want %v", x, expectedX)
	}
	if y != expectedY {
		t.Errorf("CalculateAlignedPosition() y = %v, want %v", y, expectedY)
	}
}

// ============================================================================
// Text Bounds Checking Tests
// ============================================================================

func TestCheckTextBounds(t *testing.T) {
	bbox := BBox{100, 100, 500, 500}
	textWidth, textHeight := 200.0, 50.0

	tests := []struct {
		name string
		x    float64
		y    float64
		want bool
	}{
		{"inside", 100, 150, true},
		{"exact left edge", 100, 150, true},
		{"exact right edge", 300, 150, true},
		{"exact top edge", 150, 150, true},
		{"exact bottom edge", 150, 500, true},
		{"outside left", 50, 150, false},
		{"outside right", 350, 150, false},
		{"outside top", 150, 100, false},
		{"outside bottom", 150, 550, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckTextBounds(tt.x, tt.y, textWidth, textHeight, bbox)
			if got != tt.want {
				t.Errorf("CheckTextBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTextBounds(t *testing.T) {
	x, y := 100.0, 200.0
	textWidth, textHeight := 150.0, 40.0

	want := BBox{100, 160, 250, 200}
	got := GetTextBounds(x, y, textWidth, textHeight)

	if got != want {
		t.Errorf("GetTextBounds() = %v, want %v", got, want)
	}
}

// ============================================================================
// Stroke Processing Tests
// ============================================================================

func TestGenerateStrokeOffsets(t *testing.T) {
	tests := []struct {
		name        string
		strokeWidth float64
		wantCount   int
	}{
		{"zero width", 0, 4},
		{"small width", 0.5, 4},
		{"width 1", 1, 8},
		{"large width", 5, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateStrokeOffsets(tt.strokeWidth)
			if len(got) != tt.wantCount {
				t.Errorf("GenerateStrokeOffsets() count = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

func TestGetStrokeColor(t *testing.T) {
	textColor := color.RGBA{255, 0, 0, 255}
	strokeColor := color.RGBA{0, 0, 255, 255}

	tests := []struct {
		name        string
		strokeWidth float64
		strokeColor color.Color
		textColor   color.Color
		want        color.Color
	}{
		{"no stroke", 0, strokeColor, textColor, textColor},
		{"stroke with custom color", 2, strokeColor, textColor, strokeColor},
		{"stroke nil color", 2, nil, textColor, textColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStrokeColor(tt.strokeWidth, tt.strokeColor, tt.textColor)
			if !colorsEqual(got, tt.want) {
				t.Errorf("GetStrokeColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Configuration Validation Tests
// ============================================================================

func TestValidateTextConfig(t *testing.T) {
	validBBox := BBox{0, 0, 100, 100}

	tests := []struct {
		name     string
		content  string
		fontSize float64
		bbox     BBox
		wantErr  bool
	}{
		{"valid config", "text", 24, validBBox, false},
		{"empty content", "", 24, validBBox, true},
		{"zero font size", "text", 0, validBBox, true},
		{"negative font size", "text", -10, validBBox, true},
		{"invalid bbox", "text", 24, BBox{100, 0, 0, 100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTextConfig(tt.content, tt.fontSize, tt.bbox)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTextConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// String Utilities Tests
// ============================================================================

func TestNormalizeCommand(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"center", "middle"},
		{"middle", "middle"},
		{"left", "left"},
		{"right", "right"},
		{"top", "top"},
		{"bottom", "bottom"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := NormalizeCommand(tt.input); got != tt.want {
				t.Errorf("NormalizeCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCoordinate(t *testing.T) {
	tests := []struct {
		name     string
		coordStr string
		rootSize float64
		want     float64
		wantErr  bool
	}{
		{"px value", "100px", 16, 100, false},
		{"rem value", "2rem", 16, 32, false},
		{"percentage", "50%", 16, 50, false},
		{"plain number", "100", 16, 100, false},
		{"negative number", "-50", 16, -50, false},
		{"invalid", "abc", 16, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCoordinate(tt.coordStr, tt.rootSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCoordinate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseCoordinate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func colorsEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkParseSizeValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseSizeValue("24px")
	}
}

func BenchmarkParseColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseColor("#FF0000")
	}
}

func BenchmarkParseBBox(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseBBox("10,20,100,200", BBoxXYXY)
	}
}

func BenchmarkCalculateFontSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateFontSize(2, UnitREM, 16, 100)
	}
}

func BenchmarkParseAlignment(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseAlignment("bottom+10px", 16)
	}
}