package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/donbarrigon/forge/log"
)

type Environment struct {
	App    App               `json:"app"`
	Server Server            `json:"server"`
	Db     Db                `json:"db"`
	Log    Log               `json:"log"`
	Mail   Mail              `json:"mail"`
	More   map[string]string `json:"more"`
}

type App struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	URL    string `json:"url"`
	Locale string `json:"locale"`
	Debug  bool   `json:"debug"`
}

type Server struct {
	Port              string `json:"port"`
	HttpsEnabled      bool   `json:"httpsEnabled"`
	ReadTimeout       int    `json:"readTimeout"`
	ReadHeaderTimeout int    `json:"readHeaderTimeout"`
	WriteTimeout      int    `json:"writeTimeout"`
	IdleTimeout       int    `json:"idleTimeout"`
	MaxHeaderBytes    int    `json:"maxHeaderBytes"`
}

type Db struct {
	Driver           string `json:"driver"`
	Name             string `json:"name"`
	ConnectionString string `json:"connectionString"`
}

type Log struct {
	LevelFile   string `json:"levelFile"`
	LevelOutPut string `json:"levelOutPut"`
}

type Mail struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	FromName string `json:"fromName"`
	Identity string `json:"identity"`
}

var Env *Environment = &Environment{
	App: App{
		Name:   "NewApp",
		Key:    "",
		URL:    "http://localhost:3000",
		Locale: "es",
		Debug:  false,
	},
	Server: Server{
		Port:              "3000",
		HttpsEnabled:      false,
		ReadTimeout:       30,
		ReadHeaderTimeout: 30,
		WriteTimeout:      30,
		IdleTimeout:       30,
		MaxHeaderBytes:    1 << 20,
	},
	Db: Db{
		Driver:           "mongo",
		Name:             "sample_mflix",
		ConnectionString: "mongodb://localhost:27017",
	},
	Log: Log{
		LevelFile:   "debug",
		LevelOutPut: "debug",
	},
	Mail: Mail{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: "email@gmail.com",
		Password: "secreto123",
		FromName: "Don Barrigon",
		Identity: "donbarrigon@gmail.com",
	},
}

// LoadEnv busca env.json en la misma carpeta del ejecutable y lo carga sobre
// Env. json.Unmarshal solo toca las claves presentes en el archivo, así que
// cualquier campo ausente conserva el valor por defecto que ya tiene Env.
// Si el archivo no existe o el JSON es inválido, se sigue con los defaults.
func LoadEnv() {
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

	if err := json.Unmarshal(data, Env); err != nil {
		fmt.Println("Invalid env.json, falling back to defaults:", err.Error())
	}

	pushValues()
}

// pushValues empuja los valores ya cargados en Env a las variables globales
// de los paquetes que dependen de esta configuración.
func pushValues() {
	log.LV = log.ParseLevel(Env.Log.LevelFile)
	log.LVC = log.ParseLevel(Env.Log.LevelOutPut)
}
