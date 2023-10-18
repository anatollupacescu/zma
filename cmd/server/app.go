package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/anatollupacescu/zma/bmt"
	kbmt "github.com/anatollupacescu/zma/keccak256bmt"
	"github.com/labstack/echo/v4"
)

type app struct {
	collections map[string]*bmt.Bmt[[]byte]
}

type proof struct {
	Left bool   `json:"left"`
	Sum  string `json:"sum"`
}

func (a *app) proof(c echo.Context) error {
	name := c.Param("collection")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "collection name expected")
	}

	index := c.Param("index")
	i, err := strconv.Atoi(index)
	if err != nil {
		errMsg := fmt.Sprintf("invalid integer: '%s'", index)
		return c.JSON(http.StatusBadRequest, errMsg)
	}

	cbmt, found := a.collections[name]
	if !found {
		errMsg := fmt.Sprintf("collection not found: '%s'", name)
		return c.JSON(http.StatusBadRequest, errMsg)
	}

	if i >= cbmt.Len() {
		return c.JSON(http.StatusBadRequest, "index out of range")
	}

	var proofs []proof
	for _, v := range cbmt.Proof(i) {
		proofs = append(proofs, proof{
			Left: v.Left,
			Sum:  hex.EncodeToString(v.Sum),
		})
	}

	return c.JSON(http.StatusOK, proofs)
}

func (a *app) upload(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	name := c.Param("collection")
	cbmt, found := a.collections[name]
	if !found {
		cbmt = kbmt.NewKeccak256Bmt()
		a.collections[name] = cbmt
	}

	dataHash := kbmt.Sum(data)
	root := cbmt.Add(dataHash)

	err = os.WriteFile(fmt.Sprintf("%s-%d", name, cbmt.Len()-1), data, 0644)
	if err != nil {
		return fmt.Errorf("write file to disk: %w", err)
	}

	return c.String(http.StatusOK, hex.EncodeToString(root))
}

func (a *app) download(c echo.Context) error {
	name := c.Param("collection")
	index := c.Param("index")
	i, err := strconv.Atoi(index)
	if err != nil {
		errMsg := fmt.Sprintf("invalid integer: '%s'", index)
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}
	filename := fmt.Sprintf("%s-%d", name, i)
	return c.File(filename)
}
