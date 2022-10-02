package gitOps

import (
	"fmt"
	"github.com/digtux/laminar/pkg/metrics"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"io/ioutil"
	"os"

	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/shared"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func (c *Client) Pull(state *RepoState) (err error) {
	metrics.Pulls.WithLabelValues(state.repoCfg.Name)
	var repo *git.Repository
	repo, err = git.PlainOpen(state.repoCfg.GetRealPath())
	if err == nil {
		var worktree *git.Worktree
		worktree, err = repo.Worktree()
		if err == nil {
			err = worktree.Pull(&git.PullOptions{
				RemoteName: "origin",
				Depth:      1,
			})
			if err == git.NoErrAlreadyUpToDate {
				c.logger.Debugf("pull success, already up-to-date")
				err = nil
			} else {
				state.hasBeenRead = true
			}
		}
	}
	if err == nil {
		c.logger.Debugw("repository pulled", "commitId", getCommitId(state.repoCfg))
	}
	return
}

// cloneFresh All-In-One method that will do a clone and checkout
func (c *Client) cloneFresh(repoConfig *cfg.GitRepo) (err error) {
	c.logger.Debugw("cloning",
		"url", repoConfig.URL,
		"branch", repoConfig.Branch,
		"key", repoConfig.Key,
	)
	err = c.purgeIfExists(repoConfig)
	if err == nil {
		err = c.cloneAndCheckout(repoConfig)
	}
	if err != nil {
		c.logger.Errorw("failed to cloneFresh", err)
	}
	return
}

func (c *Client) purgeIfExists(repoConfig *cfg.GitRepo) (err error) {
	diskPath := repoConfig.GetRealPath()
	if shared.IsDir(diskPath, c.logger) {
		c.logger.Debugw("previous checkout detected.. purging it",
			"path", diskPath,
		)
		err = os.RemoveAll(diskPath)
	}
	if err != nil {
		c.logger.Errorw("failed to purge", err)
	}
	return err
}

func (c *Client) cloneAndCheckout(repoConfig *cfg.GitRepo) (err error) {
	metrics.Clones.WithLabelValues(repoConfig.Name)
	mergeRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", repoConfig.Branch))

	var repo *git.Repository
	repo, err = git.PlainClone(repoConfig.GetRealPath(), false, &git.CloneOptions{
		URL:      repoConfig.URL,
		Progress: nil,
		Auth: &gitHttp.BasicAuth{
			Username: "laminar",
			Password: repoConfig.Token,
		},
		SingleBranch:  true,
		NoCheckout:    false,
		ReferenceName: mergeRef,
	})
	if err == nil {
		err = repo.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{"refs/*:refs/*"},
		})
		if err == nil {
			var worktree *git.Worktree
			worktree, err = repo.Worktree()
			if err == nil {
				err = worktree.Checkout(&git.CheckoutOptions{
					Branch: mergeRef,
				})
			}
		}
	}
	return
}

// getRemoteUpdates will check for a .laminar.yaml in the top level of a git repo
// and attempt to return []Updates from there
func getRemoteUpdates(path string) (updates *cfg.RemoteUpdates, err error) {
	rawYaml, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, ".laminar.yaml"))
	if err == nil {
		updates, err = cfg.ParseUpdates(rawYaml)
	}
	return
}

func getCommitId(repoCfg *cfg.GitRepo) string {
	repo, err := git.PlainOpen(repoCfg.GetRealPath())
	if err == nil {
		ref, err := repo.Head()
		if err == nil {
			commit, err := repo.CommitObject(ref.Hash())
			if err == nil {
				return fmt.Sprint(commit.Hash)
			}
		}
	}
	return ""
}
