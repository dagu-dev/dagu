package persistence

import (
	"fmt"
	"github.com/dagu-dev/dagu/internal/dag"
	"github.com/dagu-dev/dagu/internal/grep"
	"github.com/dagu-dev/dagu/internal/persistence/model"
	"time"
)

var (
	ErrRequestIdNotFound = fmt.Errorf("request id not found")
	ErrNoStatusDataToday = fmt.Errorf("no status data today")
	ErrNoStatusData      = fmt.Errorf("no status data")
)

type (
	DataStoreFactory interface {
		NewHistoryStore() HistoryStore
		NewDAGStore() DAGStore
	}

	HistoryStore interface {
		Open(dagFile string, t time.Time, requestId string) error
		Write(st *model.Status) error
		Close() error
		Update(dagFile, requestId string, st *model.Status) error
		ReadStatusHist(dagFile string, n int) []*model.StatusFile
		ReadStatusToday(dagFile string) (*model.Status, error)
		FindByRequestId(dagFile string, requestId string) (*model.StatusFile, error)
		RemoveAll(dagFile string) error
		RemoveOld(dagFile string, retentionDays int) error
		Rename(oldDAGFile, newDAGFile string) error
	}

	DAGStore interface {
		Create(name string, tmpl []byte) (string, error)
		List() ([]dag.DAG, error)
		Grep(pattern string) (ret []*GrepResult, errs []string, err error)
		Load(name string) (*dag.DAG, error)
		MoveDAG(oldDAGPath, newDAGPath string) error
	}

	GrepResult struct {
		Name    string
		DAG     *dag.DAG
		Matches []*grep.Match
	}
)
