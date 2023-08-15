package utils

import (
	"fmt"
	"html/template"

	"github.com/gofiber/template/html/v2"
)

func Csrf(token interface{}) template.HTML {
	csrf := fmt.Sprintf("<input type=\"hidden\" name=\"csrf\" value=\"%s\">", "")
	if str, ok := token.(string); ok {
		csrf = fmt.Sprintf("<input type=\"hidden\" name=\"csrf\" value=\"%s\">", str)
	}
	return template.HTML(csrf)
}

func SetTampletatags(engine *html.Engine) *html.Engine {
	engine.AddFunc("Csrf", Csrf)
	return engine
}
