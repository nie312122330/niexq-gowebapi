package ginext

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	//注册翻译器
	zh := zh.New()
	uni = ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	//获取gin的校验器
	validate := binding.Validator.Engine().(*validator.Validate)
	//注册翻译器
	zh_translations.RegisterDefaultTranslations(validate, trans)
}

//ValidReq 抛出参数验证的错误
func ValidReq(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindJSON(obj); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			throwErr := ValidZhError{Err: "参数验证失败", ZhErr: Translate(&verrs)}
			panic(throwErr)
		} else {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				zhErr := map[string][]string{"str": {jsonErr.Error()}}
				throwErr := ValidZhError{Err: "请求参数JSON格式错误", ZhErr: zhErr}
				panic(throwErr)
			} else {
				panic(err)
			}
		}
	}
}

//Translate 翻译错误信息
func Translate(verrs *validator.ValidationErrors) map[string][]string {
	var result = make(map[string][]string)
	for _, err := range *verrs {
		result[err.Field()] = append(result[err.Field()], err.Translate(trans))
	}
	return result
}
