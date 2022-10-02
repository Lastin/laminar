package gitOps

import (
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

type GitConfig struct {
	GitUser  string
	GitEmail string
}

type Client struct {
	logger    *zap.SugaredLogger
	cacheDb   *buntdb.DB
	states    []RepoState
	gitConfig GitConfig
}

type RepoState struct {
	repoCfg     *cfg.GitRepo
	cloned      bool
	hasBeenRead bool
}

func New(logger *zap.SugaredLogger, gitConfig GitConfig, repos []cfg.GitRepo, cacheDb *buntdb.DB) *Client {
	return &Client{
		logger:    logger,
		gitConfig: gitConfig,
		cacheDb:   cacheDb,
		states:    getEmptyStates(repos),
	}
}

func getEmptyStates(repos []cfg.GitRepo) (states []RepoState) {
	states = make([]RepoState, len(repos))
	for i, repoCfg := range repos {
		states[i] = RepoState{
			repoCfg: &repoCfg,
		}
	}
	return
}

func (c *Client) UpdateAll() (err error) {
	for _, state := range c.states {
		err = c.updateState(&state)
		if err != nil {
			c.logger.Errorw("failed to update the repo state", "state.repoCfg.Name", state.repoCfg.Name)
			return
		}
	}
	return err
}

func (c *Client) updateState(state *RepoState) (err error) {
	if !state.cloned {
		err = c.cloneFresh(state.repoCfg)
	} else {
		err = c.Pull(state)
	}
	if err == nil {
		err = state.reviewRemoteUpdates()
	}

	return
}

func (c *Client) GetStates() []RepoState {
	return c.states
}

func (state *RepoState) reviewRemoteUpdates() (err error) {
	// This sections deals with loading remote config from the gitoperations repo
	// if RemoteConfig is set we want to attempt to read '.laminar.yaml' from the remote repo
	if state.repoCfg.RemoteConfig {
		var remoteUpdates *cfg.RemoteUpdates
		remoteUpdates, err = getRemoteUpdates(state.repoCfg.GetRealPath())
		if err == nil {
			state.repoCfg.SetUpdates(remoteUpdates)
		}
	}
	return
}

func (state *RepoState) GetFilesToScan() (result []string) {
	if state.hasBeenRead {
		result = state.repoCfg.GetAllFilePaths()
		state.hasBeenRead = true
	}
	return
}
