package nlp

import (
	"fmt"

	"github.com/audibailey/XenonDiscord/config"
	"github.com/audibailey/XenonDiscord/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/lexruntimeservice"
)

func ConfigureAWS() *AWS {
	sess, err :=  session.NewSession(&aws.Config{
		Region: aws.String(config.Conf.NLP.Region),
		Credentials: credentials.NewStaticCredentials(config.Conf.NLP.AKID, config.Conf.NLP.AKSEC, ""),
	})
	if err != nil {
		logger.Log.Error("Fatal error with NLP integrator: ", err)
	}
	mySession := session.Must(sess, nil)
	svc := lexruntimeservice.New(mySession)
	return &AWS{
			Lex: svc,
		}
}

func (a *AWS) Request(input string, id string) (intent string, params map[string]interface{}, err error) {
	logger.Log.Info("The request is being parsed by AWS:Lex")
	botalias := config.Conf.NLP.BotAlias
	botname := config.Conf.NLP.BotName

	data := &lexruntimeservice.PostTextInput{}
	data.SetBotAlias(botalias)
	data.SetBotName(botname)
	data.SetInputText(input)
	data.SetUserId(id)
	err = data.Validate()
	if err == nil {
		result, err := a.Lex.PostText(data)
		if err != nil {
			return "", nil, err
		}

		if *result.DialogState == "ReadyForFulfillment" {
			var lexRes = make(map[string]interface{})

			 for key, val := range result.Slots {
				 lexRes[key] = *val
			 }

			return *result.IntentName, lexRes, nil
		}
	} else {
		logger.Log.Error(err)
		return "", nil, err
	}

	return "", nil, fmt.Errorf("NLP Failed")

}
