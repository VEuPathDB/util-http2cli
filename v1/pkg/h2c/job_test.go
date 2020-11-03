package h2c_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestJob_CreateJobDir(t *testing.T) {
	Convey("Job.CreateJobDir", t, func() {
		Convey("calls configured function", func() {
			val := 0
			target := h2c.NewJob()
			target.MkDir = func(string, os.FileMode) error {
				val++
				return nil
			}
			target.Config = &h2c.Config{OutDir: "/"}

			So(target.CreateJobDir(), ShouldBeNil)
			So(val, ShouldEqual, 1)
		})

		Convey("returns an error", func() {
			Convey("when directory creation fails", func() {
				target := h2c.NewJob()
				target.Config = &h2c.Config{OutDir: "/"}
				target.MkDir = func(string, os.FileMode) error {
					return errors.New("test")
				}

				So(target.CreateJobDir(), ShouldNotBeNil)
			})
		})
	})
}

func TestJob_GetJobPath(t *testing.T) {
	Convey("Job.GetJobPath", t, func() {
		Convey("returns a path consisting of the base output directory and the job id", func() {
			config := h2c.Config{OutDir: "/out"}
			target := h2c.NewJob()
			target.Config = &config
			target.ID = "testing"

			So(target.GetJobPath(), ShouldEqual, filepath.Join(config.OutDir, target.ID))
		})
	})
}

func TestJob_GetErrorFilePath(t *testing.T) {
	Convey("Job.GetErrorFilePath", t, func() {
		Convey("returns a path consisting of the base output directory, the job id, and the configured error file name", func() {
			config := h2c.Config{OutDir: "/out"}
			target := h2c.NewJob()
			target.Config = &config
			target.ID = "testing"

			So(target.GetErrorFilePath(), ShouldEqual, filepath.Join(config.OutDir, target.ID, h2c.ErrorFileName))
		})
	})
}

func TestJob_GetLogFilePath(t *testing.T) {
	Convey("Job.GetLogFilePath", t, func() {
		Convey("returns a path consisting of the base output directory, the job id, and the configured error file name", func() {
			config := h2c.Config{OutDir: "/out"}
			target := h2c.NewJob()
			target.Config = &config
			target.ID = "testing"

			So(target.GetLogFilePath(), ShouldEqual, filepath.Join(config.OutDir, target.ID, h2c.LogFileName))
		})
	})
}

func TestJob_CreateLogs(t *testing.T) {
	Convey("Job.CreateLogs", t, func() {
		Convey("returns an error", func() {
			Convey("if creating the log file fails", func() {
				config := h2c.Config{OutDir: "/out"}
				target := h2c.NewJob()
				target.Config = &config
				target.ID = "foobar"
				target.MkFile = func(s string) (*os.File, error) {
					if s == filepath.Join(config.OutDir, target.ID, h2c.LogFileName) {
						return nil, errors.New("test error")
					}
					return nil, nil
				}

				_, _, err := target.CreateLogs()
				So(err, ShouldNotBeNil)
			})

			Convey("if creating the error log file fails", func() {
				config := h2c.Config{OutDir: "/out"}
				target := h2c.NewJob()
				target.Config = &config
				target.ID = "foobar"
				target.MkFile = func(s string) (*os.File, error) {
					if s == filepath.Join(config.OutDir, target.ID, h2c.ErrorFileName) {
						return nil, errors.New("test error")
					}
					return nil, nil
				}

				_, _, err := target.CreateLogs()
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestJob_VerifyJobDir(t *testing.T) {
	Convey("Job.VerifyJobDir", t, func() {
		Convey("returns an error", func() {
			Convey("if the directory stat errors", func() {
				config := h2c.Config{OutDir: "/out"}
				target := h2c.NewJob()
				target.Config = &config
				target.ID = "foobar"
				target.Stat = func(string) (os.FileInfo, error) {
					return nil, errors.New("some error")
				}

				So(target.VerifyJobDir(), ShouldNotBeNil)
			})

			Convey("if the directory exists", func() {
				config := h2c.Config{OutDir: "/out"}
				target := h2c.NewJob()
				target.Config = &config
				target.ID = "foobar"
				target.Stat = func(string) (os.FileInfo, error) {
					return new(fakeFileInfo), nil
				}

				So(target.VerifyJobDir(), ShouldNotBeNil)
			})
		})

		Convey("returns nil", func() {
			Convey("if the directory does not already exist", func() {
				config := h2c.Config{OutDir: "/out"}
				target := h2c.NewJob()
				target.Config = &config
				target.ID = "foobar"
				target.Stat = func(string) (os.FileInfo, error) {
					return nil, os.ErrNotExist
				}

				So(target.VerifyJobDir(), ShouldBeNil)
			})
		})
	})
}

type fakeFileInfo struct{}

func (fakeFileInfo) Name() string       { panic("don't call me") }
func (fakeFileInfo) Size() int64        { panic("don't call me") }
func (fakeFileInfo) Mode() os.FileMode  { panic("don't call me") }
func (fakeFileInfo) ModTime() time.Time { panic("don't call me") }
func (fakeFileInfo) IsDir() bool        { panic("don't call me") }
func (fakeFileInfo) Sys() interface{}   { panic("don't call me") }
