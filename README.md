# XenonDiscord
Yeah I know, no comments. Honestly this was a proof of concept that kinda just
worked. I will finalise adding each NLP provider but have no clue when I'll
actually tidy up this project. For now there is an example in cmd on how to use
it, just run with `go run XenonBot.go --config=/path/to/config` since I'm using
Viper I think you can use any config format viper supports but I just used YAML.

Config Example

``` yaml
Debug: true
Logging:
    Logger: 0 // 1 means logrus, 0 means zap
    EnableConsole: true
    ConsoleLevel: "debug"
    ConsoleJSONFormat: false
    EnableFile: true
    FileLevel: "debug"
    FileJsonFormat: true
    FileLocation: "/path/to/log/location"
Discord:
    Token: "discordtoken"
    LoadMemebersQuietly: true
    ProjectName: "botname"
    DisableCache: true
NLP:
    Service: "aws"
    AKID: "key"
    AKSEC: "seckey"
    BotAlias: "botalias"
    BotName: "botname"
    Region: "us-west-2"
```
