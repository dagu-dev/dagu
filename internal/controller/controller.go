package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/yohamta/dagman/internal/config"
	"github.com/yohamta/dagman/internal/database"
	"github.com/yohamta/dagman/internal/models"
	"github.com/yohamta/dagman/internal/scheduler"
	"github.com/yohamta/dagman/internal/sock"
	"github.com/yohamta/dagman/internal/utils"
)

type Controller interface {
	Stop() error
	Start(bin string, workDir string, params string) error
	Retry(bin string, workDir string, reqId string) error
	GetStatus() (*models.Status, error)
	GetLastStatus() (*models.Status, error)
	GetStatusByRequestId(requestId string) (*models.Status, error)
	GetStatusHist(n int) ([]*models.StatusFile, error)
	UpdateStatus(*models.Status) error
}

func GetDAGs(dir string) (dags []*DAG, errs []string, err error) {
	dags = []*DAG{}
	errs = []string{}
	if !utils.FileExists(dir) {
		errs = append(errs, fmt.Sprintf("invalid DAGs directory: %s", dir))
		return
	}
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("%v", err)
	}
	for _, fi := range fis {
		if filepath.Ext(fi.Name()) != ".yaml" {
			continue
		}
		dag, err := fromConfig(filepath.Join(dir, fi.Name()), true)
		if err != nil {
			log.Printf("%v", err)
			if dag == nil {
				errs = append(errs, err.Error())
				continue
			}
		}
		dags = append(dags, dag)
	}
	return dags, errs, nil
}

var _ Controller = (*controller)(nil)

type controller struct {
	cfg *config.Config
}

func New(cfg *config.Config) Controller {
	return &controller{
		cfg: cfg,
	}
}

func (c *controller) Stop() error {
	client := sock.Client{Addr: sock.GetSockAddr(c.cfg.ConfigPath)}
	_, err := client.Request("POST", "/stop")
	return err
}

func (c *controller) Start(bin string, workDir string, params string) (err error) {
	go func() {
		args := []string{"start"}
		if params != "" {
			args = append(args, fmt.Sprintf("--params=\"%s\"", params))
		}
		args = append(args, c.cfg.ConfigPath)
		cmd := exec.Command(bin, args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
		cmd.Dir = workDir
		cmd.Env = os.Environ()
		defer cmd.Wait()
		err = cmd.Start()
		if err != nil {
			log.Printf("failed to start a DAG: %v", err)
		}
	}()
	time.Sleep(time.Millisecond * 500)
	return
}

func (c *controller) Retry(bin string, workDir string, reqId string) (err error) {
	go func() {
		args := []string{"retry"}
		args = append(args, fmt.Sprintf("--req=%s", reqId))
		args = append(args, c.cfg.ConfigPath)
		cmd := exec.Command(bin, args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
		cmd.Dir = workDir
		cmd.Env = os.Environ()
		defer cmd.Wait()
		err := cmd.Start()
		if err != nil {
			log.Printf("failed to retry a DAG: %v", err)
		}
	}()
	time.Sleep(time.Millisecond * 500)
	return
}

func (s *controller) GetStatus() (*models.Status, error) {
	client := sock.Client{Addr: sock.GetSockAddr(s.cfg.ConfigPath)}
	ret, err := client.Request("GET", "/status")
	if err != nil {
		if errors.Is(err, sock.ErrTimeout) {
			return nil, err
		} else {
			return defaultStatus(s.cfg), nil
		}
	}
	status, err := models.StatusFromJson(ret)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (s *controller) GetLastStatus() (*models.Status, error) {
	client := sock.Client{Addr: sock.GetSockAddr(s.cfg.ConfigPath)}
	ret, err := client.Request("GET", "/status")
	if err == nil {
		return models.StatusFromJson(ret)
	}
	if err != nil && errors.Is(err, sock.ErrTimeout) {
		return nil, err
	}
	db := database.New(database.DefaultConfig())
	status, err := db.ReadStatusToday(s.cfg.ConfigPath)
	if err != nil {
		if err != database.ErrNoDataFile {
			fmt.Printf("read status failed : %s", err)
		}
		return defaultStatus(s.cfg), nil
	}
	return status, nil
}

func (s *controller) GetStatusByRequestId(requestId string) (*models.Status, error) {
	db := database.New(database.DefaultConfig())
	ret, err := db.FindByRequestId(s.cfg.ConfigPath, requestId)
	if err != nil {
		return nil, err
	}
	return ret.Status, nil
}

func (s *controller) GetStatusHist(n int) ([]*models.StatusFile, error) {
	db := database.New(database.DefaultConfig())
	ret, err := db.ReadStatusHist(s.cfg.ConfigPath, n)
	if err != nil {
		return []*models.StatusFile{}, nil
	}
	return ret, nil
}

func (s *controller) UpdateStatus(status *models.Status) error {
	client := sock.Client{Addr: sock.GetSockAddr(s.cfg.ConfigPath)}
	res, err := client.Request("GET", "/status")
	if err != nil {
		if errors.Is(err, sock.ErrTimeout) {
			return err
		}
	}
	if err == nil {
		ss, err := models.StatusFromJson(res)
		if err != nil {
			return err
		}
		if ss.RequestId == status.RequestId && ss.Status == scheduler.SchedulerStatus_Running {
			return fmt.Errorf("the DAG is running")
		}
	}
	db := database.New(database.DefaultConfig())
	toUpdate, err := db.FindByRequestId(s.cfg.ConfigPath, status.RequestId)
	if err != nil {
		return err
	}
	w, err := db.NewWriterFor(s.cfg.ConfigPath, toUpdate.File)
	if err != nil {
		return err
	}
	if err := w.Open(); err != nil {
		return err
	}
	defer w.Close()
	if err := w.Write(status); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func defaultStatus(cfg *config.Config) *models.Status {
	return models.NewStatus(
		cfg,
		nil,
		scheduler.SchedulerStatus_None,
		int(models.PidNotRunning), nil, nil)
}
