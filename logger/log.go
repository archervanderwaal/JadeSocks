package logger

import (
	"github.com/archervanderwaal/JadeSocks/utils"
	go_logger "github.com/phachon/go-logger"
	"os"
	"path"
)

const format = "%millisecond_format% [%level_string%] [%function%:%line%] %body%"

var Logger *go_logger.Logger

func init() {
	logDir := path.Join(utils.Home(), "logger")
	if !utils.Exists(logDir) {
		_ = os.Mkdir(logDir, 0755)
	}
	Logger = go_logger.NewLogger()
	_ = Logger.Detach("console")
	consoleConfig := &go_logger.ConsoleConfig{
		Color: true,
		JsonFormat: false,
		Format: format,
	}
	_ = Logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)
	fileConfig := &go_logger.FileConfig {
		Filename : path.Join(logDir, "all.logger"),
		LevelFileName : map[int]string {
			Logger.LoggerLevel("error"): path.Join(logDir, "error.logger"),
			Logger.LoggerLevel("info"):  path.Join(logDir, "info.logger"),
			Logger.LoggerLevel("debug"): path.Join(logDir, "debug.logger"),
		},
		MaxSize : 1024 * 1024,
		MaxLine : 100000,
		DateSlice : "d",
		JsonFormat: false,
		Format: format,
	}
	_ = Logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
}