package h2c_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/veupathdb/http2cli/v1/pkg/h2c"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewMetaController(t *testing.T) {
	Convey("NewMetaController", t, func() {
		Convey("does not return nil", func() {
			So(h2c.NewMetaController(nil), ShouldNotBeNil)
		})
	})
}

func TestMeta_ServeHTTP(t *testing.T) {
	Convey("MetaController.ServeHTTP", t, func() {
		Convey("writes the expected output to the response writer", func() {
			config := new(h2c.Config)
			config.Tools = []string{"test"}
			config.Version = "test version"

			target := h2c.NewMetaController(config)

			writer := httptest.NewRecorder()

			target.ServeHTTP(writer, nil)

			So(writer.Body, ShouldNotBeNil)

			tmp := make(map[string]interface{})
			err := json.Unmarshal(writer.Body.Bytes(), &tmp)

			So(err, ShouldBeNil)

			Convey("with the server version included", func() {
				if tmp, exists := tmp["version"]; exists {
					So(tmp, ShouldEqual, "test version")
				} else {
					So(exists, ShouldBeTrue)
				}
			})

			Convey("with the server uptime included", func() {
				if tmp, exists := tmp["uptime"]; exists {
					So(len(tmp.(string)), ShouldBeGreaterThan, 0)
					d, err := time.ParseDuration(tmp.(string))
					So(err, ShouldBeNil)
					So(d, ShouldBeGreaterThan, 0)
				} else {
					So(exists, ShouldBeTrue)
				}
			})

			if tmp, exists := tmp["tools"]; exists {
				if slc, converts := tmp.([]interface{}); converts {
					So(len(slc), ShouldEqual, 1)
					So(slc[0], ShouldEqual, "test")
				} else {
					So(converts, ShouldBeTrue)
				}
			} else {
				So(exists, ShouldBeTrue)
			}

		})
	})
}
