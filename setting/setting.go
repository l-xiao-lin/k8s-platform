package setting

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name           string `mapstructure:"name"`
	Mode           string `mapstructure:"mode"`
	Version        string `mapstructure:"version"`
	Port           int    `mapstructure:"port"`
	Kubeconfig     string `mapstructure:"kubeconfig"`
	PodLogTailLine int    `mapstructure:"pod_log_tail_line"`
	WsPort         int    `mapstructure:"ws_port"`
	StartTime      string `mapstructure:"start_time"`
	MachineID      int64  `mapstructure:"machine_id"`
	*LogConfig     `mapstructure:"log"`
	*MysqlConfig   `mapstructure:"mysql"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	DbName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

var Conf = new(AppConfig)

func Init() (err error) {
	viper.SetConfigFile("./conf/config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper ReadInConfig failed,err:%v\n", err)
		return
	}

	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper Unmarshal failed,err:%v\n", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper Unmarshal failed,err:%v\n", err)
			return
		}
	})
	return

}
