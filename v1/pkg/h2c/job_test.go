package h2c_test

import (
	"errors"
	"os"
	"os/exec"
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

func TestJob_Run(t *testing.T) {
	config := new(h2c.Config)
	config.OutDir = "/out"
	config.DbDir  = "/db"
	config.Version = "test version"

	jobID := "test id"

	Convey("Job.Run", t, func() {
		Convey("returns an error", func() {
			Convey("when job output directory verification fails", func() {
				target := h2c.NewJob()
				target.Config = config
				target.ID = jobID
				target.Stat = func(string) (os.FileInfo, error) {
					return nil, errors.New("test error 1")
				}

				So(target.Run().Error(), ShouldEqual, "test error 1")
			})

			Convey("when job output directory creation fails", func() {
				target := h2c.NewJob()
				target.Config = config
				target.ID = jobID
				target.Stat = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
				target.MkDir = func(string, os.FileMode) error {
					return errors.New("test error 2")
				}

				So(target.Run().Error(), ShouldEqual, "test error 2")
			})

			Convey("when job log files creation fails", func() {
				target := h2c.NewJob()
				target.Config = config
				target.ID = jobID
				target.Stat = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
				target.MkDir = func(string, os.FileMode) error { return nil }
				target.MkFile = func(string) (*os.File, error) {
					return nil, errors.New("test error 3")
				}

				So(target.Run().Error(), ShouldEqual, "test error 3")
			})
		})

		Convey("should call job executor on success", func() {
				testFile1 := &os.File{}
				testFile2 := &os.File{}
				calls := 0

				target := h2c.NewJob()
				target.Tool = "test tool"
				target.Args = []string{"test1", "test2", "test3"}
				target.Config = config
				target.ID = jobID
				target.Stat = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
				target.MkDir = func(string, os.FileMode) error { return nil }
				target.MkFile = func(string) (out *os.File, err error) {
					if calls == 0 {
						out = testFile1
					} else {
						out = testFile2
					}

					calls++

					return
				}
				target.Exec = func(cmd *exec.Cmd) {
					So(cmd, ShouldNotBeNil)
					So(cmd.Args[0], ShouldEqual, target.Tool)
					So(cmd.Args[1:], ShouldResemble, target.Args)
					So(cmd.Stdout, ShouldPointTo, testFile1)
					So(cmd.Stderr, ShouldPointTo, testFile2)
				}

				So(target.Run(), ShouldEqual, nil)
				So(calls, ShouldEqual, 2)
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
