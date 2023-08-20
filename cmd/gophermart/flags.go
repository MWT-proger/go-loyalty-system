package main

import (
	"flag"

	"github.com/MWT-proger/go-loyalty-system/configs"
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags(conf *configs.Config) {

	flag.StringVar(&conf.HostServer, "a", conf.HostServer, "адрес и порт для запуска сервера")
	flag.StringVar(&conf.DatabaseDSN, "d", conf.DatabaseDSN, "строка с адресом подключения к БД")
	flag.StringVar(&conf.LogLevel, "l", "info", "уровень логирования")
	flag.Parse()
}
