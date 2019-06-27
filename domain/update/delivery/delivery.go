package delivery

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"github.com/bloom42/rz-go/v2"
	"github.com/bloom42/rz-go/v2/log"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/labstack/echo"
	"github.com/zhuharev/whu/domain/update"
	"github.com/zhuharev/whu/domain/webhook"
)

func New(e *echo.Echo, updUC update.UseCase, whRepo webhook.Repo) {
	s := &srv{updUC, whRepo}
	e.GET("/webhooks/:xid/updates", s.handleUpdates)
	e.POST("/webhooks/:xid", s.handleWH)
	e.POST("/create", s.handleWHCreate)
}

type srv struct {
	repo   update.UseCase
	whRepo webhook.Repo
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

func handleSlackChallenge(ctx echo.Context, body []byte) error {
	var req struct {
		Challenge string
	}
	err := json.Unmarshal(body, &req)
	if err != nil {
		return errors.Wrap(err, "unmarshal json")
	}
	if req.Challenge != "" {
		return ctx.JSON(200, map[string]string{
			"challenge": req.Challenge,
		})
	}
	return nil
}

func (s *srv) handleWH(ctx echo.Context) (err error) {
	log.Debug("incoming webhook", rz.String("url", ctx.Request().URL.String()))
	//TODO: check wh is exists
	var data []byte
	defer ctx.Request().Body.Close()
	if ct := ctx.Request().Header.Get("Content-Type"); strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		var values url.Values
		values, err = ctx.FormParams()
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

	// check is slack challenge request
	// TODO: do not handle every request
	if ct := ctx.Request().Header.Get("Content-Type"); strings.HasPrefix(ct, "application/json") {
		err = handleSlackChallenge(ctx, data)
		if err != nil {
			log.Error("handle slack challenge", rz.Err(err))
			return err
		}
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

func (s *srv) handleWHCreate(ctx echo.Context) error {
	id := xid.New()
	err := s.whRepo.Create(id.String())
	if err != nil {
		return err
	}
	return ctx.JSON(200, id.String())
}
