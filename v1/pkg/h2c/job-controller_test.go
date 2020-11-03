package h2c_test

import (
	"bytes"
	"context"
	"github.com/gorilla/mux"
	"github.com/veupathdb/http2cli/v1/pkg/h2c"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewJobController(t *testing.T) {
	Convey("NewJobController", t, func() {
		Convey("does not return nil", func() {
			So(h2c.NewJobController(nil), ShouldNotBeNil)
		})
	})
}

func TestEndpoint_ServeHTTP(t *testing.T) {
	Convey("JobController.ServeHTTP", t, func() {
		Convey("returns a 400 error if the body is empty", func() {
			req := new(http.Request)
			res := httptest.NewRecorder()

			target := h2c.NewJobController(nil)
			target.ServeHTTP(res, req)

			So(res.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("returns a 404 error if the given tool name is not found", func() {
			conf := new(h2c.Config)
			conf.Tools = []string{"grep"}

			req, _ := http.NewRequestWithContext(
				context.Background(),
				http.MethodPost,
				"",
				fakeReadCloser{bytes.NewBufferString("hola")},
			)

			res := httptest.NewRecorder()

			target := h2c.NewJobController(conf)
			target.ServeHTTP(res, req)

			So(res.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("returns a 400 error if the body cannot be parsed as json", func() {
			conf := new(h2c.Config)
			conf.Tools = []string{"grep"}

			req := mux.SetURLVars(new(http.Request), map[string]string{
				"tool": "grep",
			})
			req.Body = fakeReadCloser{bytes.NewBufferString("hola")}

			res := httptest.NewRecorder()

			target := h2c.NewJobController(conf)
			target.ServeHTTP(res, req)

			So(res.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

type fakeReadCloser struct {
	buffer *bytes.Buffer
}

func (f fakeReadCloser) Read(p []byte) (n int, err error) {
	return f.buffer.Read(p)
}

func (fakeReadCloser) Close() error {
	return nil
}
