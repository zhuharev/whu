package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/zhuharev/whu/domain/update/delivery"
	"github.com/zhuharev/whu/domain/update/repo"
)

func main() {
	repo, err := repo.New(os.Getenv("DB_PATH"))
	if err != nil {
		panic(err)
	}
	e := echo.New()
	delivery.New(e, repo)
	e.Start(":" + os.Getenv("PORT"))
}
