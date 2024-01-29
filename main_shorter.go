package gas

import (
	"fmt"

	"github.com/spf13/viper"
	appd "appdynamics"
)

type ControllerConfig struct {
	Host      string `mapstructure:"host"`
	Port      string    `mapstructure:"port"`
	UseSSL    bool   `mapstructure:"useSSL"`
	Account   string `mapstructure:"account"`
	AccessKey string `mapstructure:"accessKey"`
}

type AppContextConfig struct {
	AppName  string `mapstructure:"appName"`
	TierName string `mapstructure:"tierName"`
	NodeName string `mapstructure:"nodeName"`
}

type MiscConfig struct {
	InitTimeoutMs int `mapstructure:"initTimeoutMs"`
}

type AppConfig struct {
	Controller  ControllerConfig  `mapstructure:"controller"`
	AppContext  AppContextConfig  `mapstructure:"appContext"`
	Misc        MiscConfig        `mapstructure:"misc"`
}


func MainConfig() {
	config := viper.New()
	config.SetConfigFile("appd_config.yaml")
	config.AddConfigPath(".")

	err := config.ReadInConfig()
	if err != nil {
		fmt.Errorf("MainConfig function failed to read config file: %s, error: %s", "appd_config.yaml", err)
	}

	
	// Unmarshal the configuration values into the AppConfig struct
	var ucfg AppConfig
	err = config.Unmarshal(&ucfg)
	if err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
		return
	}	

	// Configure AppD
	cfg := appd.Config{}

	// Controller
	cfg.Controller.Host = ucfg.Controller.Host
	cfg.Controller.Port = ucfg.Controller.Port
	cfg.Controller.UseSSL = ucfg.Controller.UseSSL
	cfg.Controller.Account = ucfg.Controller.Account
	cfg.Controller.AccessKey = ucfg.Controller.AccessKey

	// App Context
	cfg.AppName = ucfg.AppContext.AppName
	cfg.TierName = ucfg.AppContext.TierName
	cfg.NodeName = ucfg.AppContext.NodeName

	// misc
	cfg.InitTimeoutMs = ucfg.Misc.InitTimeoutMs

	// init the SDK - Only for Linux
	if err := appd.InitSDK(&cfg); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}
}