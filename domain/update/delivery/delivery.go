package delivery

import (
	"io/ioutil"
	"strconv"

	"github.com/bloom42/rz-go/v2"
	"github.com/bloom42/rz-go/v2/log"

	"github.com/labstack/echo"
	"github.com/zhuharev/whu/domain/update"
)

func New(e *echo.Echo, repo update.Repo) {
	s := &srv{repo}
	e.GET("/:xid/updates", s.handleUpdates)
	e.POST("/:xid", s.handleWH)
}

type srv struct {
	repo update.Repo
}

func (s *srv) handleWH(ctx echo.Context) error {
	//TODO: check wh is exists
	defer ctx.Request().Body.Close()
	data, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	return s.repo.Save(ctx.Param("xid"), data)
}

func (s *srv) handleUpdates(ctx echo.Context) error {
	xid := ctx.Param("xid")
	offset, _ := strconv.ParseInt(ctx.QueryParam("offset"), 10, 64)
	updates, err := s.repo.Get(xid, int(offset))
	if err != nil {
		log.Error("err get updates", rz.Error("err", err))
		return err
	}

	if updates == nil { // return [] instead null
		return ctx.JSON(200, []update.Update{})
	}
	return ctx.JSON(200, updates)
}
