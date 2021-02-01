package hcolor

import "fmt"

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

const (
	TextBlack   = iota + 30 // 黑色
	TextRed                 // 红色
	TextGreen               // 绿色
	TextYellow              // 黄色
	TextBlue                // 蓝色
	TextMagenta             // 紫红色
	TextCyan                // 青蓝色
	TextWhite               // 白色
)

// Black 黑色
func Black(msg string) string {
	return SetColor(msg, 0, 0, TextBlack)
}

// Red 红色
func Red(msg string) string {
	return SetColor(msg, 0, 0, TextRed)
}

// Green 绿色
func Green(msg string) string {
	return SetColor(msg, 0, 0, TextGreen)
}

// Yellow 黄色
func Yellow(msg string) string {
	return SetColor(msg, 0, 0, TextYellow)
}

// Blue 蓝色
func Blue(msg string) string {
	return SetColor(msg, 0, 0, TextBlue)
}

// Magenta 紫红色
func Magenta(msg string) string {
	return SetColor(msg, 0, 0, TextMagenta)
}

// Cyan 青蓝色
func Cyan(msg string) string {
	return SetColor(msg, 0, 0, TextCyan)
}

// White 白色
func White(msg string) string {
	return SetColor(msg, 0, 0, TextWhite)
}

// SetColor 设置屏幕打印信息的颜色
//  @msg  		打印内容
//  @conf 		配置、终端默认设置
//  @bg   		背景色、终端默认设置
//  @textColor 	前景色
func SetColor(msg string, conf, bg, textColor int) string {
	// 0x1B是标记，配置，背景颜色，前景颜色，内容，0代表恢复默认颜色
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, textColor, msg, 0x1B)
}
