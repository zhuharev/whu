package delivery

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

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

func dedupeValues(val url.Values) map[string]string {
	res := make(map[string]string)
	for key, vals := range val {
		if len(vals) == 1 {
			res[key] = vals[0]
		} else {
			log.Error("error dedupe url values", rz.String("key", key), rz.String("value", strings.Join(vals, "|")))
		}
	}
	return res
}

func (s *srv) handleWH(ctx echo.Context) (err error) {
	//TODO: check wh is exists
	var data []byte
	defer ctx.Request().Body.Close()
	if ct := ctx.Request().Header.Get("Content-Type"); strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		values, err := ctx.FormParams()
		if err != nil {
			log.Error("error get form params", rz.Error("err", err))
			return err
		}
		data, err = json.Marshal(dedupeValues(values))
	} else {
		data, err = ioutil.ReadAll(ctx.Request().Body)
	}
	if err != nil {
		log.Error("error create body", rz.Error("err", err))
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
