// 解析配置文件，写入到全局变量中

package setting

import (
    "log"
    "time"

    "github.com/go-ini/ini"
)

type App struct {
    PageSize        int
    JwtSecret       string
    RuntimeRootPath string

    ImagePrefixUrl string
    ImageSavePath  string
    ImageMaxSize   int
    ImageAllowExts []string
}

var AppSetting App

type Server struct {
    RunMode      string
    HttpPort     int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

var ServerSetting Server

type Database struct {
    Type     string
    User     string
    Password string
    Host     string
    DbName   string
}

var DatabaseSetting Database

type Log struct {
    LogLeval    string
    LogSavePath string
    LogSaveName string
    LogFileExt  string
    TimeFormat  string
}

var LogSetting Log

// 读取配置文件，填充配置项的结构体
func Setup() {
    Cfg, err := ini.Load("gin-blog/conf/app.ini")
    if err != nil {
        log.Fatalf("Fail to parse 'gin-blog/conf/app.ini': %v", err)
    }

    err = Cfg.Section("app").MapTo(&AppSetting)
    if err != nil {
        log.Fatalf("Cfg.MapTo AppSetting error: %v", err)
    }
    AppSetting.ImageMaxSize *= (1024 * 1024) // MB -> Byte

    err = Cfg.Section("server").MapTo(&ServerSetting)
    if err != nil {
        log.Fatalf("Cfg.MapTo ServerSetting error: %v", err)
    }
    ServerSetting.ReadTimeout *= time.Second
    ServerSetting.WriteTimeout *= time.Second

    err = Cfg.Section("database").MapTo(&DatabaseSetting)
    if err != nil {
        log.Fatalf("Cfg.MapTo DatabaseSetting error: %v", err)
    }

    err = Cfg.Section("log").MapTo(&LogSetting)
    if err != nil {
        log.Fatalf("Cfg.MapTo LogSetting error: %v", err)
    }
}
