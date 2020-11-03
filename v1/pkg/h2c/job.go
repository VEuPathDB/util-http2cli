package h2c

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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
)

func NewJob() *Job {
	return &Job{
		MkDir:  os.MkdirAll,
		MkFile: os.Create,
		Stat:   os.Stat,
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

// VerifyJobDir confirms that the job directory didn't already exist before this
// run.
func (j *Job) VerifyJobDir() error {
	_, err := j.Stat(j.GetJobPath())

	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}

		return err
	}

	return errors.New("a job with id " + j.ID + " already exists.")
}

// Run executes the configured job asynchronously.
//
// The job will send back a signal on completion of either `nil` if the job
// completed successfully, or an error instance if the job failed.
//
// If this function returns an error immediately on execution, that means the
// job setup failed before the job could be started.
func (j *Job) Run(pipe chan<- error) error {
	if err := j.VerifyJobDir(); err != nil {
		return err
	}

	sout, serr, err := j.CreateLogs()
	if err != nil {
		return err
	}
	defer func() {
		_ = sout.Close()
		_ = serr.Close()
	}()

	cmd := exec.Command(j.Tool, j.Args...)
	cmd.Stdout = sout
	cmd.Stderr = serr

	go func() {
		if err := cmd.Start(); err != nil {
			pipe <- err
		}

		pipe <- cmd.Wait()
	}()

	return nil
}
