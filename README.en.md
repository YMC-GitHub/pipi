# Pipi - Watermark Text Tool

**Draw watermark text on images** with support for multiple font size units, alignment options, stroke effects, and font formats (including TTC font collections), designed for flexible watermarking needs.

## ✨ Core Features
1. **`text`** Draw watermark text on images
2. **`help`** Display help information
3. **`version`** Display version information

## 📌 Complete Parameter Reference

### text Command Parameters

| Parameter | Description | Default | Example |
|-----------|-------------|---------|---------|
| --content | Text content (required) | - | "Hello World" |
| --size | Font size, supports px, %, rem, bh | 24px | 24px, 5%, 2rem, 0.1bh |
| --root-size | Root font size for rem unit | 16px | 16px |
| --color | Text color, supports color names or hex | white | red, #FF0000, rgba(255,0,0,255) |
| --stroke | Stroke width | 0 | 3 |
| --stroke-width | Alias for --stroke | 0 | 3 |
| --stroke-color | Stroke color | Same as text color | black, #000000 |
| --bbox | Bounding box coordinates | Entire image | 100,100,500,500 |
| --bbox-type | Bbox type: xyxy, xywh, cxcy | xyxy | xywh |
| --align-y | Y-axis alignment: baseline, top, bottom, middle | baseline | bottom+10px, middle-5% |
| --align-x | X-axis alignment: left, right, middle | left | right-20px, middle+10% |
| --font | Font file path or system font name | Built-in English font | /path/to/font.ttf, Microsoft YaHei |
| --input | Input image path | Creates transparent image | input.png |
| --imgz | Image size when no input provided | 1280x720 | 1920x1080 |
| --output | Output image path (required) | - | output.png |
| --debug | Enable debug output | false | - |
| -h, --help | Show help information | - | - |
| -v, --version | Show version information | - | - |

## 🚀 Most Common Command Examples

### 1. Basic Usage (Recommended)
```bash
# Simple white text
pipi text --content "Hello" --output out.png

# Specify font size and color
pipi text --content "Watermark" --size 48px --color red --output out.png
```

### 2. Using Different Units
```bash
# CSS standard rem unit (relative to root font size)
pipi text --content "Hello" --size 2rem --root-size 16 --output out.png

# Percentage of bounding box height
pipi text --content "Hello" --size 10% --output out.png

# Bounding box height multiplier
pipi text --content "Hello" --size 0.1bh --output out.png
```

### 3. Stroke Text
```bash
# Black stroke
pipi text --content "STROKE" --color white --stroke 3 --stroke-color black --output out.png

# Using alias parameter
pipi text --content "Stroke Text" --stroke-width 2 --output out.png
```

### 4. Font Support (TTC Fonts)
```bash
# Using Microsoft YaHei font
pipi text --content "Microsoft YaHei Test" --font /fonts/msyh.ttc --output out.png

# Using system font
pipi text --content "Chinese Watermark" --font "Microsoft YaHei" --size 36px --output out.png
```

### 5. Alignment and Offset
```bash
# Bottom alignment with 10px upward offset
pipi text --content "Bottom Text" --align-y bottom-10px --output out.png

# Horizontal and vertical center
pipi text --content "Center Text" --align-x middle --align-y middle --output out.png

# Relative offset (percentage)
pipi text --content "Offset Text" --align-x middle+10% --align-y bottom-5% --output out.png
```

### 6. Bounding Box Control
```bash
# XYXY format (top-left and bottom-right coordinates)
pipi text --content "Box" --bbox "100,100,500,300" --output out.png

# XYWH format (top-left coordinates + width/height)
pipi text --content "Box" --bbox "100,100,400,200" --bbox-type xywh --output out.png

# CXCY format (center coordinates + width/height)
pipi text --content "Box" --bbox "300,200,400,200" --bbox-type cxcy --output out.png
```

### 7. Custom Image Size
```bash
# Create 1920x1080 transparent background without input image
pipi text --content "Watermark" --imgz 1920x1080 --output watermark.png
```

### 8. Debug Mode
```bash
# View detailed execution information
pipi text --content "Debug" --debug --output out.png
```

## 📁 Supported Font Formats

The tool supports multiple font formats and loading methods:

| Type | Format | Description |
|------|--------|-------------|
| **Font File** | .ttf, .ttc, .otf | Load font file directly |
| **System Font** | Font name | Search system fonts via fc-list |
| **Built-in Font** | GoRegular | English only |

### System Font Lookup Example (Linux/macOS):
```bash
# Using system font name
pipi text --content "System Font" --font "Arial" --output out.png
pipi text --content "Chinese Text" --font "Source Han Sans" --output out.png
```

## 🎯 Size Units Explained

| Unit | Description | Formula | Example |
|------|-------------|---------|---------|
| **px** | Absolute pixels | Direct value | `--size 24px` |
| **%** | Percentage of bbox height | Bbox height × percentage | `--size 10%` |
| **rem** | CSS root unit | Root font size × multiplier | `--size 2rem --root-size 16` |
| **bh** | Bbox height multiplier | Bbox height × multiplier | `--size 0.1bh` |

## 🎨 Color Format Support

| Format | Example | Description |
|--------|---------|-------------|
| **Color Name** | red, green, blue, white, black, yellow, cyan, magenta | Predefined colors |
| **Hexadecimal** | #FF0000, #FF0000FF | RGB or RGBA format |
| **RGBA Function** | rgba(255,0,0,255) | Standard RGBA format |

## 🔧 Advanced Usage

### Precise Alignment Control
```bash
# Supports pixel and rem offsets
pipi text --content "Offset" --align-x left+50px --align-y top-2rem --output out.png

# Supports percentage offsets (relative to bounding box)
pipi text --content "Relative" --align-x middle+20% --align-y middle-15% --output out.png
```

### Bounding Box Type Description
```bash
# XYXY: Directly specify bounding box range
--bbox "x1,y1,x2,y2"

# XYWH: Specify top-left corner and dimensions
--bbox "x,y,width,height" --bbox-type xywh

# CXCY: Specify center point and dimensions
--bbox "centerX,centerY,width,height" --bbox-type cxcy
```

### Debug Output Example
```bash
$ pipi text --content "Test" --debug --output test.png
Loading input image: input.png
Input image size: 1920x1080
Drawing context size: 1920x1080
Using default bbox (full image): [0, 0, 1920, 1080]
Font size in px: 24px
Text: 'Test'
Text size: 48 x 30
Text position: (0, 30)
BBox bounds: X[0-1920] Y[0-1080]
Drawing main text with color: {255 255 255 255} at (0, 30)
✅ Watermark text added successfully: test.png
```

## 🐳 Container Deployment (Docker)

### Build Image

```sh
docker build --progress=plain -f Dockerfile.pipi --target runtime -t ymc/pipi .
```

### Run Conversion Tasks (Mount Local Directories)

```sh
# Start container and mount Windows system fonts
docker run -it --rm -v "$(pwd):/app" -v "/mnt/i/capture2:/data" -v "/mnt/c/Windows/Fonts:/fonts:ro" ymc/pipi bash

# Install fonts (if needed)
# apt-get update && apt-get install -y fonts-wqy-microhei
# ls /usr/share/fonts/truetype/wqy/wqy-microhei.ttc

# --font /data/fonts/simhei.ttf
# --font /usr/share/fonts/truetype/wqy/wqy-microhei.ttc
# fonts-wqy-zenhei

# fc-list :lang=zh
# WenQuanYi Zen Hei
pipi text --content "KaiTi Test" --output /data/result.png --input /data/bg_0323.png --font "/fonts/simkai.ttf"
```

## 📊 Command Quick Reference

| Command | Description | Required Parameters |
|---------|-------------|---------------------|
| text | Draw watermark text | --content, --output |
| help | Display help information | None |
| version | Display version information | None |

## ⚠️ Important Notes

1. **Output path** (`--output`) and **text content** (`--content`) are required parameters
2. **Chinese/Unicode support** requires specifying a font file or system font that contains the characters
3. **System font lookup** depends on the `fc-list` command (Linux/macOS). Windows users need to specify font file paths directly
4. **TTC font collections** automatically select the first font. Use font extraction tools if you need a specific index
5. **Bounding box warning**: Debug mode will detect if text exceeds bounding box boundaries
6. **Stroke performance**: Large stroke widths may impact performance; use reasonable stroke widths
7. **Color transparency**: Supports RGBA format for semi-transparent effects (e.g., `rgba(255,0,0,128)`)

## 🐳 Environment Requirements

- **Go Version**: 1.16+
- **Supported Platforms**: Linux, macOS, Windows (Git Bash or WSL recommended)
- **System Dependencies**: Linux/macOS require `fontconfig` (fc-list command) for system font lookup

### Install Font Tools (Linux):
```bash
# Ubuntu/Debian
sudo apt-get install fontconfig

# CentOS/RHEL
sudo yum install fontconfig
```

## 🛡️ Feature Highlights

- ✅ **Multiple unit support**: px, %, rem, bh for various sizing needs
- ✅ **TTC font support**: Full support for TrueType Collection font files
- ✅ **Precise alignment**: Supports pixel, rem, and percentage offsets
- ✅ **Smart bounding box**: Supports multiple bbox formats and auto-detection
- ✅ **Stroke effects**: 8-direction stroke algorithm for text readability
- ✅ **Debug mode**: Detailed execution information for troubleshooting
- ✅ **Transparent background**: Create transparent backgrounds without input images

## 📝 Frequently Asked Questions

**Q: Chinese characters display as boxes?**  
A: Specify a font file that contains Chinese characters, e.g., `--font "Microsoft YaHei"` or `--font /path/to/chinese-font.ttf`

**Q: How to use system fonts on Windows?**  
A: Windows doesn't support fc-list. Specify the font file path directly, e.g., `--font C:\Windows\Fonts\msyh.ttc`

**Q: How to choose font size units?**  
A: 
- Use `px` for absolute sizes
- Use `%` or `bh` for sizes relative to bounding box
- Use `rem` for CSS standard compliance

**Q: How to create semi-transparent text?**  
A: Use RGBA color format: `--color rgba(255,255,255,128)`

**Q: What if text exceeds the bounding box?**  
A: 
- Reduce font size
- Use `--debug` to view detailed information
- Adjust alignment or offset values

**Q: Does it support dynamic watermarks?**  
A: Currently supports static text only. Consider using scripts for batch generation

**Q: How to generate multiple watermarks?**  
A: Call the command multiple times or write a loop script
