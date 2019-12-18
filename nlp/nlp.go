package nlp

import (
	"github.com/audibailey/XenonDiscord/config"
	_ "github.com/audibailey/XenonDiscord/logger"

	"github.com/aws/aws-sdk-go/service/lexruntimeservice"
)

var Nlp NLP = (*NaturalLangProc)(nil)

type NaturalLangProc struct {
	*AWS
}

type AWS struct {
	Lex *lexruntimeservice.LexRuntimeService
} 

type NLP interface {
	Request(input string, id string) (intent string, params map[string]interface{}, err error)
}

func Configure() {
	if config.Conf.NLP.Service == "aws" {
		Nlp = ConfigureAWS()
	}
}  

func (*NaturalLangProc) Request(input string, id string) (intent string, params map[string]interface{}, err error) {
	return Nlp.Request(input, id)
}
