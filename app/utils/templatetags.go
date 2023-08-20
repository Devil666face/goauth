package utils

import (
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var (
	Ctx *fiber.Ctx
)

func SetTampletatags(engine *html.Engine) *html.Engine {
	engine.AddFunc("Csrf", Csrf)
	engine.AddFunc("UrlTo", UrlTo)
	engine.AddFunc("Url", Url)
	return engine
}

func Csrf(token interface{}) template.HTML {
	csrf := fmt.Sprintf("<input type=\"hidden\" name=\"csrf\" value=\"%s\">", "")
	if str, ok := token.(string); ok {
		csrf = fmt.Sprintf("<input type=\"hidden\" name=\"csrf\" value=\"%s\">", str)
	}
	return template.HTML(csrf)
}

func UrlTo(name string, key string, param interface{}) string {
	return getRouteUrl(name, fiber.Map{key: param})
}

func Url(name string) string {
	return getRouteUrl(name, fiber.Map{})
}

func getRouteUrl(name string, fmap fiber.Map) string {
	url, _ := Ctx.GetRouteURL(name, fmap)
	if url == "" {
		panic(fmt.Sprintf("Not found url for \"%s\" with params %s", name, fmap))
	}
	return url
}
