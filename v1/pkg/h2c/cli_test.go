package h2c_test

import (
	"os"
	"testing"

	"github.com/veupathdb/http2cli/v1/pkg/h2c"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInitCLI(t *testing.T) {
	Convey("InitCLI", t, func() {
		Convey("sets the config tools based on the CLI input", func() {
			os.Args = []string{"foo", "-tbar", "-tfizz", "-tbuzz"}

			tmp := new(h2c.Config)

			h2c.InitCLI(tmp)

			So(tmp.Tools, ShouldResemble, []string{"bar", "fizz", "buzz"})
			So(tmp.DbDir, ShouldEqual, h2c.DefaultDbDir)
			So(tmp.OutDir, ShouldEqual, h2c.DefaultOutDir)
			So(tmp.ServerPort, ShouldEqual, h2c.DefaultPort)
		})

		Convey("sets the config port based on the CLI input", func() {
			os.Args = []string{"foo", "-p90"}

			tmp := new(h2c.Config)

			h2c.InitCLI(tmp)

			So(tmp.Tools, ShouldResemble, h2c.DefaultTools)
			So(tmp.DbDir, ShouldEqual, h2c.DefaultDbDir)
			So(tmp.OutDir, ShouldEqual, h2c.DefaultOutDir)
			So(tmp.ServerPort, ShouldEqual, 90)
		})

		Convey("sets the config output dir on the CLI input", func() {
			os.Args = []string{"foo", "-o/bar"}

			tmp := new(h2c.Config)

			h2c.InitCLI(tmp)

			So(tmp.Tools, ShouldResemble, h2c.DefaultTools)
			So(tmp.DbDir, ShouldEqual, h2c.DefaultDbDir)
			So(tmp.OutDir, ShouldEqual, "/bar")
			So(tmp.ServerPort, ShouldEqual, h2c.DefaultPort)
		})

		Convey("sets the config db dir on the CLI input", func() {
			os.Args = []string{"foo", "-d/bar"}

			tmp := new(h2c.Config)

			h2c.InitCLI(tmp)

			So(tmp.Tools, ShouldResemble, h2c.DefaultTools)
			So(tmp.DbDir, ShouldEqual, "/bar")
			So(tmp.OutDir, ShouldEqual, h2c.DefaultOutDir)
			So(tmp.ServerPort, ShouldEqual, h2c.DefaultPort)
		})
	})
}
