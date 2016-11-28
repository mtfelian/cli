package cli

import (
	"fmt"
	"strings"
)

// список возможных кодов цветов
const (
	Black = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// список возможных SGR-кодов
const (
	Reset = iota
	Bold
	Faint
	Italic
	Underline
)

const (
	// ClearScreen это escape-код для очистки экрана
	ClearScreen = "\x1b[2J"
	// Move это escape-код для перемещения курсора на x, y
	Move = "\x1b[%d;%dH"
	// LineReset это escape-код для возврата курсора в начало строки и очистки её
	LineReset = "\r\x1b[K"
)

// getColor возвращает escape-код для цвета с кодом code
func getColor(code int) string {
	return getParam(30 + code)
}

// getBgColor возвращает escape-код для фона цвета с кодом code
func getBgColor(code int) string {
	return getParam(40 + code)
}

// getParam возвращает escape-код для параметров текста
func getParam(code int) string {
	return fmt.Sprintf("\x1b[%dm", code)
}

// Println печатает строку str и добавляет в конец перенос строки
func Println(str string, a ...interface{}) {
	Printf(str+"\n", a...)
}

// Printf печатает строку str и подставляет туда параметры
func Printf(str string, a ...interface{}) {
	fmt.Printf(Colorize(str), a...)
}

// Colorize возвращает цветную строку преобразуя {-теги
// Пример: cli.Colorize("{Rred string{0 and {Bblue part{0")
func Colorize(str string, a ...interface{}) string {
	str = fmt.Sprintf(str, a...)
	changeMap := map[string]string{
		"{w":  getColor(White),
		"{a":  getColor(Black),
		"{r":  getColor(Red),
		"{g":  getColor(Green),
		"{y":  getColor(Yellow),
		"{b":  getColor(Blue),
		"{m":  getColor(Magenta),
		"{c":  getColor(Cyan),
		"{W":  getParam(Bold) + getColor(White),
		"{A":  getParam(Bold) + getColor(Black),
		"{R":  getParam(Bold) + getColor(Red),
		"{G":  getParam(Bold) + getColor(Green),
		"{Y":  getParam(Bold) + getColor(Yellow),
		"{B":  getParam(Bold) + getColor(Blue),
		"{M":  getParam(Bold) + getColor(Magenta),
		"{C":  getParam(Bold) + getColor(Cyan),
		"{_w": getBgColor(White),
		"{_a": getBgColor(Black),
		"{_r": getBgColor(Red),
		"{_g": getBgColor(Green),
		"{_y": getBgColor(Yellow),
		"{_b": getBgColor(Blue),
		"{_m": getBgColor(Magenta),
		"{_c": getBgColor(Cyan),
		"{_W": getParam(Bold) + getBgColor(White),
		"{_A": getParam(Bold) + getBgColor(Black),
		"{_R": getParam(Bold) + getBgColor(Red),
		"{_G": getParam(Bold) + getBgColor(Green),
		"{_Y": getParam(Bold) + getBgColor(Yellow),
		"{_B": getParam(Bold) + getBgColor(Blue),
		"{_M": getParam(Bold) + getBgColor(Magenta),
		"{_C": getParam(Bold) + getBgColor(Cyan),
		"{i":  getParam(Italic),
		"{u":  getParam(Underline),
		"{0":  getParam(Reset),
		"{s":  ClearScreen,
	}
	for key, value := range changeMap {
		str = strings.Replace(str, key, value, -1)
	}
	return str
}
