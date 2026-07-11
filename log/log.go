package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LogLevel int

type Logger struct {
	Time    string   `json:"time,omitempty"`
	Level   LogLevel `json:"level"`
	Message string   `json:"message"`
	Line    string   `json:"line,omitempty"`
	File    string   `json:"file,omitempty"`
}

const (
	LV_EMERGENCY LogLevel = iota // 0 - El sistema está inutilizable
	LV_ALERT                     // 1 - Se necesita acción inmediata
	LV_CRITICAL                  // 2 - Fallo crítico del sistema
	LV_ERROR                     // 3 - Errores de ejecución
	LV_WARNING                   // 4 - Algo inesperado pasó
	LV_NOTICE                    // 5 - Eventos normales, pero significativos
	LV_INFO                      // 6 - Información general
	LV_DEBUG                     // 7 - Información detallada para depuración
	LV_PRINT                     // 8 - Solo imprime en consola
	LV_OFF                       // 9 - Desactiva todos los logs
)

const (
	LOG_PATH    = "./tmp/logs"
	DATE_FORMAT = "2006-01-02 15:04:05.000"
	RESET_COLOR = "\033[0m"
)

// LV es el nivel a partir del cual se guarda en archivo.
var LV LogLevel = LV_DEBUG

// LVC es el nivel a partir del cual se muestra en consola (independiente de LV).
var LVC LogLevel = LV_DEBUG

// Days es cuántos días de archivos .log se conservan. Channel diario fijo.
var Days = 30

func (lv LogLevel) String() string {
	switch lv {
	case LV_EMERGENCY:
		return "EMERGENCY"
	case LV_ALERT:
		return "ALERT"
	case LV_CRITICAL:
		return "CRITICAL"
	case LV_ERROR:
		return "ERROR"
	case LV_WARNING:
		return "WARNING"
	case LV_NOTICE:
		return "NOTICE"
	case LV_INFO:
		return "INFO"
	case LV_DEBUG:
		return "DEBUG"
	case LV_PRINT:
		return "PRINT"
	case LV_OFF:
		return "OFF"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel convierte un string (ej. "warning", "DEBUG") a su LogLevel
// correspondiente. Si no coincide con ninguno, retorna LV_DEBUG por defecto.
func ParseLevel(s string) LogLevel {
	switch strings.ToUpper(s) {
	case "EMERGENCY":
		return LV_EMERGENCY
	case "ALERT":
		return LV_ALERT
	case "CRITICAL":
		return LV_CRITICAL
	case "ERROR":
		return LV_ERROR
	case "WARNING":
		return LV_WARNING
	case "NOTICE":
		return LV_NOTICE
	case "INFO":
		return LV_INFO
	case "DEBUG":
		return LV_DEBUG
	case "PRINT":
		return LV_PRINT
	case "OFF":
		return LV_OFF
	default:
		return LV_DEBUG
	}
}

func (lv LogLevel) Color() string {
	switch lv {
	case LV_EMERGENCY:
		return "\033[91m" // rojo brillante
	case LV_ALERT:
		return "\033[95m" // magenta
	case LV_CRITICAL:
		return "\033[35m" // fucsia
	case LV_ERROR:
		return "\033[31m" // rojo
	case LV_WARNING:
		return "\033[33m" // amarillo
	case LV_NOTICE:
		return "\033[92m" // verde claro
	case LV_INFO:
		return "\033[34m" // azul
	case LV_DEBUG:
		return "\033[36m" // cian
	case LV_PRINT:
		return "\033[90m" // gris claro
	default:
		return RESET_COLOR
	}
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func Emergency(msg string, a ...any) { Log(LV_EMERGENCY, msg, a...) }
func Alert(msg string, a ...any)     { Log(LV_ALERT, msg, a...) }
func Critical(msg string, a ...any)  { Log(LV_CRITICAL, msg, a...) }
func Error(msg string, a ...any)     { Log(LV_ERROR, msg, a...) }
func Warning(msg string, a ...any)   { Log(LV_WARNING, msg, a...) }
func Notice(msg string, a ...any)    { Log(LV_NOTICE, msg, a...) }
func Info(msg string, a ...any)      { Log(LV_INFO, msg, a...) }
func Debug(msg string, a ...any)     { Log(LV_DEBUG, msg, a...) }
func Print(msg string, a ...any)     { Log(LV_PRINT, msg, a...) }

func Dump(a any) {
	fmt.Println(formatDump(a))
}

func DumpMany(vars ...any) {
	sep := strings.Repeat("-", 30)
	for i, v := range vars {
		if i > 0 {
			fmt.Println(sep)
		}
		fmt.Println(formatDump(v))
	}
}

// Log es el único punto de entrada: filtra por nivel (consola y archivo
// por separado) y despacha solo lo que corresponda, sin capas de configuración.
func Log(level LogLevel, msg string, a ...any) {
	toConsole := level <= LVC
	toFile := level != LV_PRINT && level <= LV

	if !toConsole && !toFile {
		return
	}

	_, file, line, _ := runtime.Caller(2)

	l := &Logger{
		Time:    time.Now().Format(DATE_FORMAT),
		Level:   level,
		Message: fmt.Sprintf(msg, a...),
		File:    filepath.Base(file),
		Line:    strconv.Itoa(line),
	}

	if toConsole {
		l.outputConsole()
	}

	if toFile {
		l.outputFile()
	}
}

func (l *Logger) outputConsole() {
	color := l.Level.Color()
	fmt.Printf("%s [%s%s%s] %s%s%s (%s:%s)\n",
		l.Time, color, l.Level.String(), RESET_COLOR,
		color, l.Message, RESET_COLOR, l.File, l.Line)
}

func (l *Logger) outputFile() {
	data, err := json.Marshal(l)
	if err != nil {
		fmt.Println("Log serialization error:", err.Error())
		return
	}
	writeToFile(data)
}

var (
	fileMu      sync.Mutex
	currentFile *os.File
	startOnce   sync.Once
	done        chan struct{}
)

// Start abre el archivo del día actual y lanza el goroutine que
// duerme hasta la medianoche siguiente para rotar. Se dispara una sola vez,
// Cuando se carga la configuracion.
func Start() {
	startOnce.Do(func() {
		done = make(chan struct{})
		rotateFile()
		go rotationLoop()
	})
}

func rotationLoop() {
	for {
		select {
		case <-time.After(time.Until(nextMidnight())):
			rotateFile()
		case <-done:
			return
		}
	}
}

func nextMidnight() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, now.Location())
}

// rotateFile abre el archivo del día actual y dispara la limpieza de
// archivos viejos. Solo se llama una vez por día (al iniciar y a cada
// medianoche), nunca por cada log.
func rotateFile() {
	today := time.Now().Format("2006-01-02")

	if err := os.MkdirAll(LOG_PATH, os.ModePerm); err != nil {
		fmt.Println("Failed to create log directory:", err.Error())
		return
	}

	f, err := os.OpenFile(filepath.Join(LOG_PATH, today+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err.Error())
		return
	}

	fileMu.Lock()
	old := currentFile
	currentFile = f
	fileMu.Unlock()

	if old != nil {
		old.Close()
	}

	go deleteOldFiles(today)
}

func writeToFile(data []byte) {
	defer fileMu.Unlock()
	fileMu.Lock()
	if currentFile == nil {
		// Se libera antes de intentar abrir: rotateFile toma fileMu por su cuenta
		fileMu.Unlock()
		rotateFile()
		fileMu.Lock()
		if currentFile == nil {
			fmt.Println("Log file is not initialized.")
			return
		}
	}

	currentFile.Write(data)
	currentFile.Write([]byte("\n"))
}

// deleteOldFiles corre una vez por rotación de día, no en cada log.
func deleteOldFiles(today string) {
	if Days <= 0 {
		return
	}

	cutoff, err := time.Parse("2006-01-02", today)
	if err != nil {
		return
	}
	cutoff = cutoff.AddDate(0, 0, -Days)

	entries, err := os.ReadDir(LOG_PATH)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
			continue
		}
		datePart := strings.TrimSuffix(entry.Name(), ".log")
		entryDate, err := time.Parse("2006-01-02", datePart)
		if err == nil && entryDate.Before(cutoff) {
			_ = os.Remove(filepath.Join(LOG_PATH, entry.Name()))
		}
	}
}

// Close cierra el archivo de log actual y detiene el goroutine de rotación.
// hay que llamarlo al apagar la aplicación.
func Close() {
	if done != nil {
		close(done)
	}

	fileMu.Lock()
	defer fileMu.Unlock()
	if currentFile != nil {
		currentFile.Close()
		currentFile = nil
	}
}
