package config

import (
	"flag"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Debug bool `mapstructure:"Debug"`
	Logging loggingConfig `mapstructure:"Logging"`
	Discord discordConfig `mapstructure:"Discord"`
	NLP nlpConfig `mapstructure:"NLP"`
	OWM string `mapstructure:"OWM"`
}

type loggingConfig struct {
	Logger int `mapstructure:"Logger"`
	EnableConsole bool `mapstructure:"EnableConsole"`
	ConsoleLevel string `mapstructure:"ConsoleLevel"`
	ConsoleJSONFormat bool `mapstructure:"ConsoleJSONFormat"`
	EnableFile bool `mapstructure:"EnableFile"`
	FileLevel string `mapstructure:"FileLevel"`
	FileJSONFormat bool `mapstructure:"FileJSONFormat"`
	FileLocation string `mapstructure:"FileLocation"`
}

type discordConfig struct {
	Token string `mapstructure:"Token"`
	LoadMembersQuietly bool `mapstructure:"LoadMembersQuietly"`
	ProjectName string `mapstructure:"ProjectName"`
	DisableCache bool `mapstructure:"DisableCache"`
}

type nlpConfig struct {
	Service string `mapstructure:"Service"`
	AKID string `mapstructure:"AKID"`
	AKSEC string `mapstructure:"AKSEC"`
	BotAlias string `mapstructure:"BotAlias"`
	BotName string `mapstructure:"BotName"`
	Region string `mapstructure:"Region"`
}

var (
	CfgFile string
	Conf *Config
)


func Configure(configdata *Config){
	if configdata == nil {
		flag.StringVar(&CfgFile, "config", "", "config file")

		pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
		pflag.Parse()
		viper.BindPFlags(pflag.CommandLine)

		CfgFile, err := homedir.Expand(viper.GetString("config"))
		if err != nil {
			log.Fatal("Could not expand config location: ", err.Error())
		}

		if CfgFile != "" {
			viper.SetConfigFile(CfgFile)
		} else {
			log.Fatal("You must have a config file flag!")
		}

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Fatal("Could not find config file: ", err.Error())
			} else {
				log.Fatal("Error with config file: ", err.Error())
			}
		} else {
			Conf = &Config{}
			err = viper.Unmarshal(Conf)
			if err != nil {
				log.Fatal("Unable to decode config: ", err.Error())
			}
		}
	} else {
		Conf = configdata
	} 
}
