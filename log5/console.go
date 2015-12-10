package log5
import (
	"log"
	"os"
)

// 设置颜色刷
type Brush func(string) string

func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []Brush{
	NewBrush("1;36"), // Trace      cyan
	NewBrush("1;34"), // Info		blue
	NewBrush("1;33"), // Warning    yellow
	NewBrush("1;31"), // Error      red
	NewBrush("1;37"), // Fatal		white

}

type ConsoleLog struct {
	log *log.Logger
	level Level
}

// 初始化控制台输出引擎
func NewConsole() LogEngine {
	return &ConsoleLog{
		log:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		level: Trace,
	}
}

func (c *ConsoleLog) Init() error {
	return nil
}

func (c *ConsoleLog) Write(msg string, level Level) error {
	if level > c.level {
		return nil
	}
	c.log.Println(colors[level](msg))

	return nil
}

func (c *ConsoleLog) Destroy() {}

func init() {
	Register("console", NewConsole)
}
