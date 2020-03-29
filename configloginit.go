package configloginit

import (
  "os"
  "strings"
  "github.com/sirupsen/logrus"
  "sync"
  "github.com/spf13/viper"
  "fmt"
)

var onceLogger sync.Once
var logger *logrus.Logger

type Formatter struct {
  TimestampFormat string
  LogFormat string
}

func InitConfig(configname string) {
    index := strings.LastIndex(configname, "/")
    if index != -1 {
        cIndex := strings.LastIndex(configname, ".yaml")
        if cIndex == -1 {
            panic(fmt.Errorf("config file not specified"))
        }
        viper.SetConfigName(configname[index:cIndex])
        viper.AddConfigPath(configname[:index])
        err := viper.ReadInConfig()
        if err != nil {
            panic(fmt.Errorf("Fatal error config file: %s \n", err))
        }
    } else {
        cIndex := strings.LastIndex(configname, ".yaml")
        if cIndex == -1 {
            panic(fmt.Errorf("config file not specified"))
        }
        fmt.Println(cIndex)
        viper.SetConfigName(configname[:cIndex])
        viper.AddConfigPath(".")
        err := viper.ReadInConfig()
        if err != nil {
            panic(fmt.Errorf("Fatal error config file: %s \n", err))
        }
    }
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

  output := f.LogFormat
  tsf := f.TimestampFormat

  output = strings.Replace(output, "%time%", entry.Time.Format(tsf), 1)
  output = strings.Replace(output, "%msg%", entry.Message, 1)
  level := strings.ToUpper(entry.Level.String())
  output = strings.Replace(output, "%lvl%", level,1)
  return []byte(output), nil
}

func InitLogger(filename string) *logrus.Logger {
    onceLogger.Do(func() {

      f, err := os.OpenFile(filename, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
      logger = &logrus.Logger{
        Level: logrus.DebugLevel,
        Formatter: &Formatter {
          TimestampFormat: "2006-01-02 15:04:05",
          LogFormat: "%time% - [%lvl%] - %msg%\n",
        },
      }

      if err != nil {
          panic(err)
      }
      logger.SetOutput(f)

    })
    return logger
}
