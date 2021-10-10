package context

import (
	"encoding/json"
	"io"
	"net/http"
)

type Conetxt struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Conetxt) ReadJson(req interface{}) error {
	r := c.R
	byte, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(byte, req)
}

func (c *Conetxt) WriteJson(code int, resp interface{}) error {
	c.W.WriteHeader(code)
	respJson, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(respJson)
	return err
}

func (c *Conetxt) OkJson(resp interface{}) error {
	return c.WriteJson(http.StatusOK, resp)
}

func (c *Conetxt) SystemErrorJson(resp interface{}) error {
	return c.WriteJson(http.StatusInternalServerError, resp)
}

func (c *Conetxt) BadRequestJson(resp interface{}) error {
	return c.WriteJson(http.StatusBadRequest, resp)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Conetxt {
	return &Conetxt{
		W: w,
		R: r,
	}
}
