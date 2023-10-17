package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/anatollupacescu/zma/bmt"
	"github.com/labstack/echo/v4"
)

type app struct {
	collections map[string]*bmt.Bmt
}

type proof struct {
	Left bool   `json:"left"`
	Sum  string `json:"sum"`
}

func (a *app) proof(c echo.Context) error {
	name := c.Param("collection")
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "collection name expected")
	}

	index := c.Param("index")
	i, err := strconv.Atoi(index)
	if err != nil {
		errMsg := fmt.Sprintf("invalid integer: '%s'", index)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	cbmt, found := a.collections[name]
	if !found {
		errMsg := fmt.Sprintf("collection not found: '%s'", name)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	var proofs []proof
	for _, v := range cbmt.Proof(i) {
		proofs = append(proofs, proof{
			Left: v.Left,
			Sum:  hex.EncodeToString(v.Sum[:]),
		})
	}

	return c.JSON(http.StatusOK, proofs)
}

func (a *app) upload(c echo.Context) error {
	name := c.Param("collection")
	c.Logger().Debug(name)

	req := c.Request()

	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	cbmt, found := a.collections[name]
	if !found {
		cbmt = bmt.New()
		a.collections[name] = cbmt
	}

	root := cbmt.Add(data)

	return c.String(http.StatusOK, hex.EncodeToString(root[:]))
}
