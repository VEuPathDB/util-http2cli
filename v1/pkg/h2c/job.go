package h2c

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	LogFileName   = "log.txt"
	ErrorFileName = "error.txt"

	DirPerms = 0755
)

type (
	DirCreateFn  func(string, os.FileMode) error
	FileCreateFn func(string) (*os.File, error)
	FileStatFn   func(string) (os.FileInfo, error)
	ExecFn       func(*exec.Cmd)
)

func NewJob() *Job {
	return &Job{
		MkDir:  os.MkdirAll,
		MkFile: os.Create,
		Stat:   os.Stat,
		Exec:   execute,
	}
}

type Job struct {
	ID     string
	Tool   string
	Args   []string
	Config *Config

	MkDir  DirCreateFn
	MkFile FileCreateFn
	Stat   FileStatFn
	Exec   ExecFn
}

// GetJobPath returns the path to the output directory for the current job.
func (j *Job) GetJobPath() string {
	return filepath.Join(j.Config.OutDir, j.ID)
}

// GetLogFilePath returns the path to the log file for the current job.
func (j *Job) GetLogFilePath() string {
	return filepath.Join(j.GetJobPath(), LogFileName)
}

// GetErrorFilePath returns the path to the error log file for the current job.
func (j *Job) GetErrorFilePath() string {
	return filepath.Join(j.GetJobPath(), ErrorFileName)
}

// CreateJobDir creates the output directory for this job.
func (j *Job) CreateJobDir() error {
	return j.MkDir(j.GetJobPath(), DirPerms)
}

// CreateLogs creates the output log files for this job.
func (j *Job) CreateLogs() (sout, serr *os.File, err error) {
	if sout, err = j.MkFile(j.GetLogFilePath()); err != nil {
		return
	}

	serr, err = j.MkFile(j.GetErrorFilePath())

	return
}

func (j *Job) VerifyJobDir() error {
	_, err := j.Stat(j.GetJobPath())

	if err != nil {
		if os.IsNotExist(err) {
			return j.MkDir(j.GetJobPath(), 0775)
		}

		return err
	}

	return nil
}

// Run executes the configured job asynchronously.
//
// The job will send back a signal on completion of either `nil` if the job
// completed successfully, or an error instance if the job failed.
//
// If this function returns an error immediately on execution, that means the
// job setup failed before the job could be started.
func (j *Job) Run() error {
	if err := j.VerifyJobDir(); err != nil {
		return err
	}

	sout, serr, err := j.CreateLogs()
	if err != nil {
		logrus.Errorf("Failed to create log files for job %s", j.ID)
		return err
	}
	defer func() {
		_ = sout.Close()
		_ = serr.Close()
	}()

	cmd := exec.Command(j.Tool, j.Args...)
	cmd.Stdout = sout
	cmd.Stderr = serr
	cmd.Dir = j.GetJobPath()

	j.Exec(cmd)

	return nil
}

func execute(cmd *exec.Cmd) {
	logrus.Debug("Executing command ", cmd)

	if err := cmd.Run(); err != nil {
		logrus.Debug("Command run failed with ", err.Error())
		_, _ = cmd.Stderr.Write([]byte(err.Error()))
	}
}