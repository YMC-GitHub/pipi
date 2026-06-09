根据您提供的文档模板和Go代码，我为您生成了中文文档 `README.zh.md`：

```markdown
# Pipi - 水印文字工具

**在图像上绘制水印文字**，支持多种字体大小单位、对齐方式、描边效果和字体格式（包括TTC字体集合），专为灵活的水印添加需求设计。

## ✨ 核心功能
1. **`text`** 在图像上绘制水印文字
2. **`help`** 显示帮助信息
3. **`version`** 显示版本信息

## 📌 完整参数说明

### text 命令参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| --content | 文本内容（必需） | - | "Hello World" |
| --size | 字体大小，支持px、%、rem、bh | 24px | 24px, 5%, 2rem, 0.1bh |
| --root-size | rem单位的根字体大小 | 16px | 16px |
| --color | 文本颜色，支持颜色名或十六进制 | white | red, #FF0000, rgba(255,0,0,255) |
| --stroke | 描边宽度 | 0 | 3 |
| --stroke-width | --stroke的别名 | 0 | 3 |
| --stroke-color | 描边颜色 | 同文本颜色 | black, #000000 |
| --bbox | 边界框坐标 | 整个图像 | 100,100,500,500 |
| --bbox-type | 边界框类型：xyxy, xywh, cxcy | xyxy | xywh |
| --align-y | Y轴对齐：baseline, top, bottom, middle | baseline | bottom+10px, middle-5% |
| --align-x | X轴对齐：left, right, middle | left | right-20px, middle+10% |
| --font | 字体文件路径或系统字体名 | 内置英文字体 | /path/to/font.ttf, Microsoft YaHei |
| --input | 输入图像路径 | 创建透明图像 | input.png |
| --imgz | 无输入图像时的图像尺寸 | 1280x720 | 1920x1080 |
| --output | 输出图像路径（必需） | - | output.png |
| --debug | 启用调试输出 | false | - |
| -h, --help | 显示帮助信息 | - | - |
| -v, --version | 显示版本信息 | - | - |

## 🚀 最常用命令示例

### 1. 基础用法（推荐）
```bash
# 简单的白色文字
pipi text --content "Hello" --output out.png

# 指定字体大小和颜色
pipi text --content "水印文字" --size 48px --color red --output out.png
```

### 2. 使用不同单位
```bash
# CSS标准rem单位（相对于根字体大小）
pipi text --content "Hello" --size 2rem --root-size 16 --output out.png

# 边界框高度百分比
pipi text --content "Hello" --size 10% --output out.png

# 边界框高度倍数
pipi text --content "Hello" --size 0.1bh --output out.png
```

### 3. 描边文字
```bash
# 黑色描边
pipi text --content "STROKE" --color white --stroke 3 --stroke-color black --output out.png

# 使用别名参数
pipi text --content "描边文字" --stroke-width 2 --output out.png
```

### 4. 中文支持（TTC字体）
```bash
# 使用微软雅黑字体
pipi text --content "微软雅黑测试" --font /fonts/msyh.ttc --output out.png

# 使用系统字体
pipi text --content "中文水印" --font "Microsoft YaHei" --size 36px --output out.png
```

### 5. 对齐和偏移
```bash
# 底部对齐并向上偏移10像素
pipi text --content "Bottom Text" --align-y bottom-10px --output out.png

# 水平居中，垂直居中
pipi text --content "Center Text" --align-x middle --align-y middle --output out.png

# 相对偏移（百分比）
pipi text --content "Offset Text" --align-x middle+10% --align-y bottom-5% --output out.png
```

### 6. 边界框控制
```bash
# XYXY格式（左上角和右下角坐标）
pipi text --content "Box" --bbox "100,100,500,300" --output out.png

# XYWH格式（左上角坐标+宽高）
pipi text --content "Box" --bbox "100,100,400,200" --bbox-type xywh --output out.png

# CXCY格式（中心点坐标+宽高）
pipi text --content "Box" --bbox "300,200,400,200" --bbox-type cxcy --output out.png
```

### 7. 自定义图像尺寸
```bash
# 不提供输入图像，创建1920x1080透明背景
pipi text --content "Watermark" --imgz 1920x1080 --output watermark.png
```

### 8. 调试模式
```bash
# 查看详细的执行信息
pipi text --content "Debug" --debug --output out.png
```

## 📁 支持的字体格式

工具支持多种字体格式和加载方式：

| 类型 | 格式 | 说明 |
|------|------|------|
| **字体文件** | .ttf, .ttc, .otf | 直接加载字体文件 |
| **系统字体** | 字体名称 | 通过fc-list查找系统字体 |
| **内置字体** | GoRegular | 仅支持英文 |

### 系统字体查找示例（Linux/macOS）：
```bash
# 使用系统字体名称
pipi text --content "System Font" --font "Arial" --output out.png
pipi text --content "思源黑体" --font "Source Han Sans" --output out.png
```

## 🎯 尺寸单位详解

| 单位 | 说明 | 计算公式 | 示例 |
|------|------|----------|------|
| **px** | 绝对像素 | 直接使用数值 | `--size 24px` |
| **%** | 边界框高度百分比 | 边界框高度 × 百分比 | `--size 10%` |
| **rem** | CSS根单位 | 根字体大小 × 倍数 | `--size 2rem --root-size 16` |
| **bh** | 边界框高度倍数 | 边界框高度 × 倍数 | `--size 0.1bh` |

## 🎨 颜色格式支持

| 格式 | 示例 | 说明 |
|------|------|------|
| **颜色名称** | red, green, blue, white, black, yellow, cyan, magenta | 预定义颜色 |
| **十六进制** | #FF0000, #FF0000FF | RGB或RGBA格式 |
| **RGBA函数** | rgba(255,0,0,255) | 标准RGBA格式 |

## 🔧 高级用法

### 精确对齐控制
```bash
# 支持像素和rem偏移
pipi text --content "Offset" --align-x left+50px --align-y top-2rem --output out.png

# 支持百分比偏移（相对于边界框）
pipi text --content "Relative" --align-x middle+20% --align-y middle-15% --output out.png
```

### 边界框类型说明
```bash
# XYXY：直接指定边界框范围
--bbox "x1,y1,x2,y2"

# XYWH：指定左上角和尺寸
--bbox "x,y,width,height" --bbox-type xywh

# CXCY：指定中心点和尺寸
--bbox "centerX,centerY,width,height" --bbox-type cxcy
```

### 调试输出示例
```bash
$ pipi text --content "Test" --debug --output test.png
```

## 🐳 容器化部署（Docker）

### 构建镜像

```sh
docker build --progress=plain -f Dockerfile.pipi --target runtime -t ymc/pipi .
```

### 运行转换任务（挂载本地目录）

```sh
# 启动容器 并 挂载 window系统字体
docker run -it --rm -v "$(pwd):/app" -v "/mnt/i/capture2:/data" -v "/mnt/c/Windows/Fonts:/fonts:ro" ymc/pipi bash

# 安装字体
# apt-get update && apt-get install -y fonts-wqy-microhei
# ls /usr/share/fonts/truetype/wqy/wqy-microhei.ttc

# --font /data/fonts/simhei.ttf
# --font /usr/share/fonts/truetype/wqy/wqy-microhei.ttc
# fonts-wqy-zenhei

# fc-list :lang=zh
# WenQuanYi Zen Hei
pipi text --content "楷体测试" --output /data/result.png --input /data/bg_0323.png --font "/fonts/simkai.ttf"

```


## 📊 命令速查表

| 命令 | 说明 | 必需参数 |
|------|------|----------|
| text | 绘制水印文字 | --content, --output |
| help | 显示帮助信息 | 无 |
| version | 显示版本信息 | 无 |

## ⚠️ 注意事项

1. **输出路径** (`--output`) 和 **文本内容** (`--content`) 是必需参数
2. **中文支持**需要指定包含中文字符的字体文件或系统字体
3. **系统字体查找**依赖 `fc-list` 命令（Linux/macOS），Windows用户需直接指定字体文件路径
4. **TTC字体集合**会自动选择第一个字体，如需指定索引请使用字体文件提取工具
5. **边界框警告**：调试模式下会检测文本是否超出边界框范围
6. **描边性能**：大尺寸描边可能会影响性能，建议使用合理的描边宽度
7. **颜色透明度**：支持RGBA格式，可实现半透明效果（如 `rgba(255,0,0,128)`）

## 🐳 环境要求

- **Go版本**：1.16+
- **支持平台**：Linux, macOS, Windows（推荐使用Git Bash或WSL）
- **系统依赖**：Linux/macOS需要 `fontconfig` (fc-list命令) 用于系统字体查找

### 安装字体工具（Linux）：
```bash
# Ubuntu/Debian
sudo apt-get install fontconfig

# CentOS/RHEL
sudo yum install fontconfig
```

## 🛡️ 特性亮点

- ✅ **多单位支持**：px、%、rem、bh，满足各种尺寸需求
- ✅ **TTC字体支持**：完整支持TrueType Collection字体集合
- ✅ **精确对齐**：支持像素、rem和百分比偏移
- ✅ **智能边界框**：支持多种边界框格式和自动检测
- ✅ **描边效果**：8方向描边算法，确保文字可读性
- ✅ **调试模式**：详细输出执行信息，便于问题排查
- ✅ **透明背景**：支持无输入图像时创建透明背景
