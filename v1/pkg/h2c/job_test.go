package h2c_test

import (
	"testing"

	"github.com/veupathdb/http2cli/v1/pkg/h2c"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewJob(t *testing.T) {
	Convey("NewJob", t, func() {
		Convey("does not return nil", func() {
			So(h2c.NewJob(), ShouldNotBeNil)
		})
	})
}
