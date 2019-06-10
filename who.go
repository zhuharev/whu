package main

import (
	"os"

	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/zhuharev/whu/domain/update"
	"github.com/zhuharev/whu/domain/update/delivery"
	"github.com/zhuharev/whu/domain/update/repo"
	whStormDB "github.com/zhuharev/whu/domain/webhook/repo/storm"
)

func main() {
	sdb, err := storm.Open(os.Getenv("DB_PATH"))
	defer sdb.Close()
	repo, err := repo.New(sdb)
	if err != nil {
		panic(err)
	}
	whRepo := whStormDB.New(sdb)
	updUC := update.New(whRepo, repo)
	e := echo.New()
	delivery.New(e, updUC, whRepo)
	e.Start(":" + os.Getenv("PORT"))
}
