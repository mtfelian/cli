package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// LineReset возвращает курсор в начало строки и очищает её
const LineReset = "\r\x1b[K"

// Список возможных цветов
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

// список возможных SGR-параметров
const (
	Reset = iota
	Bold
	Faint
	Italic
	Underline
)

// Output это указатель на буфер STDOUT
var Output *bufio.Writer = bufio.NewWriter(os.Stdout)

// getColor возвращает escape-код для цвета с кодом code
func getColor(code int) string {
	return fmt.Sprintf("\x1b[3%dm", code)
}

// getBgColor возвращает escape-код для фона цвета с кодом code
func getBgColor(code int) string {
	return fmt.Sprintf("\x1b[4%dm", code)
}

// getParam возвращает escape-код для параметров текста
func getParam(code int) string {
	return fmt.Sprintf("\x1b[%dm", code)
}

/*
Установить процентный флаг: num | PCT
Проверить процентный флаг: num & PCT
Сбросить процентный флаг: num & 0xFF
*/
const shift = uint(^uint(0)>>63) << 4
const PCT = 0x8000 << shift

// Screen это глобальный буфер экрана
var Screen *bytes.Buffer = new(bytes.Buffer)

// getXY получает относительные и абсолютные координаты
// что бы получить относительные координаты, установите флаг PCT в число.
// Пример. Получить 10% от полной ширины по x и 20 по y
// x, y = cli.GetXY(10|cli.PCT, 20)
func getXY(x int, y int) (int, int) {
	if y == -1 {
		y = CurrentHeight() + 1
	}

	if x&PCT != 0 {
		x = int((x & 0xFF) * GetWidth() / 100)
	}

	if y&PCT != 0 {
		y = int((y & 0xFF) * GetHeight() / 100)
	}

	return x, y
}

type sf func(int, string) string

// applyTransform применяет заданную функцию преобразования sf к каждой строке внутри str
func applyTransform(str string, transform sf) (out string) {
	out = ""

	for idx, line := range strings.Split(str, "\n") {
		out += transform(idx, line)
	}

	return
}

// Clear очищает экран
func Clear() {
	Output.WriteString("\x1b[2J")
}

// MoveСursor перемещает курсор в положение, заданное координатами (x, y)
func MoveCursor(x int, y int) {
	fmt.Fprintf(Screen, "\x1b[%d;%dH", x, y)
}

// MoveTo перемещает строку str в положение, заданное координатами (x, y)
func MoveTo(str string, x int, y int) (out string) {
	x, y = getXY(x, y)

	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\x1b[%d;%dH%s", y+idx, x, line)
	})
}

// ResetLine возвращает каретку в начало строки str
func ResetLine(str string) (out string) {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf(LineReset, line)
	})
}

// Color применяет к строке str заданный цвет color
// cli.Color("Red string", cli.Red)
func Color(str string, color int) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", getColor(color), line, Reset)
	})
}

// MakeBold делает строку str жирной
func MakeBold(str string) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[1m%s\033[0m", line)
	})
}

// Highlight выделяет цветом color подстроку substr в строке str
func Highlight(str, substr string, color int) string {
	hiSubstr := Color(substr, color)
	return strings.Replace(str, substr, hiSubstr, -1)
}

// HighlightRegion выделяет цветом color символы с индексами от from до to в строке str
func HighlightRegion(str string, from, to, color int) string {
	return str[:from] + Color(str[from:to], color) + str[to:]
}

// Background изменяет цвет фона строки str на color
// cli.Background("string", cli.Red)
func Background(str string, color int) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", getBgColor(color), line, Reset)
	})
}

// getWinsize получает ширину и высоту терминала
// возвращает: 1. ширину (x), 2. высоту (y), 3. ошибку/nil
func getWinsize() (int, int, error) {
	x, y := 0, 0

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return x, y, err
	}

	outputParts := strings.Split(string(out), " ")
	if len(outputParts) != 2 {
		return x, y, fmt.Errorf("Wrong output: %s", out)
	}

	x64, err := strconv.ParseInt(outputParts[0], 10, 32)
	if err != nil {
		return x, y, err
	}

	y64, err := strconv.ParseInt(outputParts[1], 10, 32)
	if err != nil {
		return x, y, err
	}

	x, y = int(x64), int(y64)
	return x, y, nil
}

// GetWidth возвращает ширину консоли
func GetWidth() int {
	x, _, err := getWinsize()
	if err != nil {
		return -1
	}
	return x
}

// GetHeight возвращает высоту консоли
func GetHeight() int {
	_, y, err := getWinsize()
	if err != nil {
		return -1
	}
	return y
}

// CurrentHeight возвращает текущую высоту (количество строк в экранном буфере)
func CurrentHeight() int {
	return strings.Count(Screen.String(), "\n")
}

// Flush записывает в буфер экрана с учётом что бы он не переполнился
func Flush() {
	for idx, str := range strings.Split(Screen.String(), "\n") {
		height := GetHeight()
		if idx > height && height > 0 {
			return
		}

		Output.WriteString(str + "\n")
	}

	Output.Flush()
	Screen.Reset()
}

// Print пишет в буфер экрана
func Print(a ...interface{}) {
	fmt.Fprint(Screen, a...)
}

// Println пишет в буфер экрана добавляя в конце символ перевода строки
func Println(a ...interface{}) {
	fmt.Fprintln(Screen, a...)
}

// Printf пишет в буфер экрана согласно заданному формату format
func Printf(format string, a ...interface{}) {
	fmt.Fprintf(Screen, format, a...)
}

// Colorize делает строку цветной согласно {-тегам
func Colorize(str string, a ...interface{}) string {
	changeMap := map[string]string{
		"{w": getColor(White),
		"{a": getColor(Black),
		"{r": getColor(Red),
		"{g": getColor(Green),
		"{y": getColor(Yellow),
		"{b": getColor(Blue),
		"{m": getColor(Magenta),
		"{c": getColor(Cyan),
		"{W": getParam(Bold) + getColor(White),
		"{A": getParam(Bold) + getColor(Black),
		"{R": getParam(Bold) + getColor(Red),
		"{G": getParam(Bold) + getColor(Green),
		"{Y": getParam(Bold) + getColor(Yellow),
		"{B": getParam(Bold) + getColor(Blue),
		"{M": getParam(Bold) + getColor(Magenta),
		"{C": getParam(Bold) + getColor(Cyan),
		"{i": getParam(Italic),
		"{u": getParam(Underline),
		"{0": getParam(Reset),
	}
	for key, value := range changeMap {
		str = strings.Replace(str, key, value, -1)
	}
	return fmt.Sprintf(str, a...)
}
