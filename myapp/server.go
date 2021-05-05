package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Static("/static", "assets")

	e.Logger.Fatal(e.Start(":1323"))
}

kubectl run -it --rm --restart=Never bb_temp --image=busybox -n monitoring -- sh -c 'while true; do echo "Start probing pushgateway...";do wget -O - -q http://:9091/-/healthy; done'