package logger

type Conf struct {
	Model      string `toml:"model"`      // console, file
	Level      string `toml:"level"`      // debug, info, warn, error, dpanic, panic, fatal
	Filename   string `toml:"filename"`   // log file name
	MaxSize    int    `toml:"maxSize"`    // max size per log file, unit is MB
	MaxAge     int    `toml:"MaxAge"`     //	max age per log file, unit is day
	MaxBackups int    `toml:"maxBackups"` // max backups per log file
	TimeFormat string `toml:"timeFormat"` // time format
	LocalTime  bool   `toml:"localTime"`  // local time
	Compress   bool   `toml:"compress"`   // compress
}
