package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/donbarrigon/forge/logs"
)

// EnvExtra es un mapa de variables de entorno propias de la app,
// con métodos de acceso tipados que devuelven el valor cero si falla.
type EnvExtra map[string]string

// Get devuelve el valor como string, o "" si no existe.
func (self EnvExtra) Get(key string) string {
	return self[key]
}

// Exists comprueba si la clave está presente en el mapa.
func (self EnvExtra) Exists(key string) bool {
	_, ok := self[key]
	return ok
}

// GetInt devuelve el valor como int, o 0 si no existe o no es válido.
func (self EnvExtra) GetInt(key string) int {
	v, err := strconv.Atoi(self[key])
	if err != nil {
		return 0
	}
	return v
}

// GetInt8 devuelve el valor como int8, o 0 si no existe o no es válido.
func (self EnvExtra) GetInt8(key string) int8 {
	v, err := strconv.ParseInt(self[key], 10, 8)
	if err != nil {
		return 0
	}
	return int8(v)
}

// GetInt16 devuelve el valor como int16, o 0 si no existe o no es válido.
func (self EnvExtra) GetInt16(key string) int16 {
	v, err := strconv.ParseInt(self[key], 10, 16)
	if err != nil {
		return 0
	}
	return int16(v)
}

// GetInt32 devuelve el valor como int32, o 0 si no existe o no es válido.
func (self EnvExtra) GetInt32(key string) int32 {
	v, err := strconv.ParseInt(self[key], 10, 32)
	if err != nil {
		return 0
	}
	return int32(v)
}

// GetInt64 devuelve el valor como int64, o 0 si no existe o no es válido.
func (self EnvExtra) GetInt64(key string) int64 {
	v, err := strconv.ParseInt(self[key], 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// GetUint devuelve el valor como uint, o 0 si no existe o no es válido.
func (self EnvExtra) GetUint(key string) uint {
	v, err := strconv.ParseUint(self[key], 10, 0)
	if err != nil {
		return 0
	}
	return uint(v)
}

// GetUint8 devuelve el valor como uint8, o 0 si no existe o no es válido.
func (self EnvExtra) GetUint8(key string) uint8 {
	v, err := strconv.ParseUint(self[key], 10, 8)
	if err != nil {
		return 0
	}
	return uint8(v)
}

// GetUint16 devuelve el valor como uint16, o 0 si no existe o no es válido.
func (self EnvExtra) GetUint16(key string) uint16 {
	v, err := strconv.ParseUint(self[key], 10, 16)
	if err != nil {
		return 0
	}
	return uint16(v)
}

// GetUint32 devuelve el valor como uint32, o 0 si no existe o no es válido.
func (self EnvExtra) GetUint32(key string) uint32 {
	v, err := strconv.ParseUint(self[key], 10, 32)
	if err != nil {
		return 0
	}
	return uint32(v)
}

// GetUint64 devuelve el valor como uint64, o 0 si no existe o no es válido.
func (self EnvExtra) GetUint64(key string) uint64 {
	v, err := strconv.ParseUint(self[key], 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// GetFloat32 devuelve el valor como float32, o 0.0 si no existe o no es válido.
func (self EnvExtra) GetFloat32(key string) float32 {
	v, err := strconv.ParseFloat(self[key], 32)
	if err != nil {
		return 0
	}
	return float32(v)
}

// GetFloat64 devuelve el valor como float64, o 0.0 si no existe o no es válido.
func (self EnvExtra) GetFloat64(key string) float64 {
	v, err := strconv.ParseFloat(self[key], 64)
	if err != nil {
		return 0
	}
	return v
}

// GetBool devuelve el valor como bool, o false si no existe o no es válido.
// Acepta "true", "1", "t", "yes", "y" como verdadero; cualquier otro es falso.
func (self EnvExtra) GetBool(key string) bool {
	v, err := strconv.ParseBool(self[key])
	if err != nil {
		return false
	}
	return v
}

// GetDuration devuelve el valor como time.Duration, o 0 si no existe o no es válido.
// Acepta formatos de time.ParseDuration (ej: "1s", "5m", "2h").
func (self EnvExtra) GetDuration(key string) time.Duration {
	d, err := time.ParseDuration(self[key])
	if err != nil {
		return 0
	}
	return d
}

// ----------------------------------------------------------------------------
// Estructuras de configuración
// ----------------------------------------------------------------------------

type Environment struct {
	App    EnvApp    `json:"app"`
	Server EnvServer `json:"server"`
	Db     EnvDB     `json:"db"`
	Log    EnvLog    `json:"log"`
	Mail   EnvMail   `json:"mail"`
	Extra  EnvExtra  `json:"extra"`
}

type EnvApp struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	URL    string `json:"url"`
	Locale string `json:"locale"`
	Debug  bool   `json:"debug"`
}

type EnvServer struct {
	Port              string `json:"port"`
	HttpsEnabled      bool   `json:"httpsEnabled"`
	ReadTimeout       int    `json:"readTimeout"`
	ReadHeaderTimeout int    `json:"readHeaderTimeout"`
	WriteTimeout      int    `json:"writeTimeout"`
	IdleTimeout       int    `json:"idleTimeout"`
	MaxHeaderBytes    int    `json:"maxHeaderBytes"`
}

type EnvDB struct {
	Driver           string          `json:"driver"`
	Name             string          `json:"name"`
	ConnectionString string          `json:"connectionString"`
	ClientOptions    DBClientOptions `json:"clientOptions"`
}

type DBClientOptions struct {
	MaxPoolSize uint64 `json:"maxPoolSize"`
	MinPoolSize uint64 `json:"minPoolSize"`
	RetryWrites bool   `json:"retryWrites"`
	Timeout     int    `json:"timeout"`
}

type EnvLog struct {
	LevelFile   string `json:"levelFile"`
	LevelOutPut string `json:"levelOutPut"`
}

type EnvMail struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	FromName string `json:"fromName"`
	Identity string `json:"identity"`
}

// ----------------------------------------------------------------------------
// Variables globales con valores por defecto
// ----------------------------------------------------------------------------

var (
	App = EnvApp{
		Name:   "NewApp",
		Key:    "",
		URL:    "http://localhost:3000",
		Locale: "es",
		Debug:  false,
	}
	Server = EnvServer{
		Port:              "3000",
		HttpsEnabled:      false,
		ReadTimeout:       30,
		ReadHeaderTimeout: 30,
		WriteTimeout:      30,
		IdleTimeout:       30,
		MaxHeaderBytes:    1 << 20,
	}
	DB = EnvDB{
		Driver:           "mongo",
		Name:             "sample_mflix",
		ConnectionString: "mongodb://localhost:27017",
		ClientOptions: DBClientOptions{
			MaxPoolSize: 100,
			MinPoolSize: 5,
			RetryWrites: true,
			Timeout:     30,
		},
	}
	Log = EnvLog{
		LevelFile:   "debug",
		LevelOutPut: "debug",
	}
	Mail = EnvMail{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: "email@gmail.com",
		Password: "secreto123",
		FromName: "Don Barrigon",
		Identity: "donbarrigon@gmail.com",
	}

	// Extra guarda variables de entorno propias de la app que no pertenecen
	// a la librería (no tienen su propio struct tipado).
	Extra = EnvExtra{}
)

// ----------------------------------------------------------------------------
// Carga y aplicación de la configuración
// ----------------------------------------------------------------------------

// Load busca env.json en la misma carpeta del ejecutable y lo carga sobre
// los globals de este paquete. json.Unmarshal solo toca las claves presentes
// en el archivo, así que cualquier campo ausente conserva el valor por
// defecto que ya tiene cada global. Si el archivo no existe o el JSON es
// inválido, se sigue con los defaults.
func Load() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Could not resolve executable path:", err.Error())
		pushValues()
		return
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(exePath), "env.json"))
	if err != nil {
		// No existe env.json o no se pudo leer: se queda con los defaults.
		pushValues()
		return
	}

	// se arma con los defaults actuales para que Unmarshal solo pise las
	// claves presentes en el JSON y deje el resto intacto
	env := Environment{
		App:    App,
		Server: Server,
		Db:     DB,
		Log:    Log,
		Mail:   Mail,
		Extra:  Extra,
	}

	if err := json.Unmarshal(data, &env); err != nil {
		fmt.Println("Invalid env.json, falling back to defaults:", err.Error())
	}

	App = env.App
	Server = env.Server
	DB = env.Db
	Log = env.Log
	Mail = env.Mail
	Extra = env.Extra

	pushValues()
}

// pushValues empuja los valores ya cargados a las variables globales de los
// paquetes que dependen de esta configuración.
func pushValues() {
	logs.LV = logs.ParseLevel(Log.LevelFile)
	logs.LVC = logs.ParseLevel(Log.LevelOutPut)
}
