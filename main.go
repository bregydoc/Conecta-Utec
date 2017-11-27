package main

import (
	"github.com/kataras/iris"
	"./hotspot"
	"log"
)

var hot *hotspot.HotspotWrapper

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	app.RegisterView(iris.HTML("./views", ".html").Reload(true))

	app.StaticWeb("/statics", "./statics")

	app.Get("/", func(c iris.Context) {
		c.View("index.html")
	})

	app.Get("/dashboard", func(c iris.Context) {
		hot = hotspot.NewHotspotWrapper("windows")
		c.View("dashboard.html")

	})

	app.Post("/api/configure", func(c iris.Context) {
		values := c.FormValues()
		ssid := values["name"][0]
		pass := values["password"][0]
		log.Println(ssid, pass)

		hot.ConfigProperties(ssid, pass)
		log.Println(hot.ConfigCommand)
		c.JSON(iris.Map{
			"status": "ok",
			"ssid": ssid,
			"key": pass,
		})

	})

	app.Post("/api/activate", func(c iris.Context) {
		err := hot.ActivateHotspot()
		if err != nil {
			c.JSON(iris.Map{
				"status": "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(iris.Map{
			"status": "ok",
			"message": "",
		})
	})

	app.Post("/api/stop", func(c iris.Context) {
		err := hot.StopHotspot()
		if err != nil {
			c.JSON(iris.Map{
				"status": "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(iris.Map{
			"status": "ok",
			"message": "",
		})
	})


	app.Run(iris.Addr(":5400"))

}

