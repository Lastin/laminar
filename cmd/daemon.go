package cmd

import (
	"fmt"
	"github.com/digtux/laminar/pkg/metrics"
	"github.com/digtux/laminar/pkg/registryCore"
	"github.com/digtux/laminar/pkg/shared"
	"github.com/digtux/laminar/pkg/web"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
	"io/ioutil"
	"os"
	"time"

	"github.com/digtux/laminar/pkg/cache"
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/gitOps"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "daemon",
	Short: "launch laminar daemon service/server",
	Long: `This is the laminar server and API

Laminar is a GitOps utility for automating the promotion of docker images in gitoperations.

.. turbulence is bad.. use laminar flow`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := New()
		if err != nil {
			panic(err)
		}
		d.Start()
	},
}

type Daemon struct {
	logger               *zap.SugaredLogger
	dockerRegistryClient *registryCore.Client
	webClient            *web.Client
	cacheDb              *buntdb.DB
	fileList             []string // list of files containing docker images urls
	gitConfig            cfg.Global
	gitOpsClient         *gitOps.Client
}

func New() (d *Daemon, err error) {
	logger := shared.GetLogger(debug)
	appConfig, err := loadConfig(logger)
	if err != nil {
		return nil, err
	}
	cacheDb := cache.Open(configCache, logger)
	if err != nil {
		return
	}
	d = &Daemon{
		logger:               logger,
		dockerRegistryClient: registryCore.New(logger, appConfig.DockerRegistries, cacheDb),
		webClient:            web.New(logger, appConfig.Global.GitHubToken),
		cacheDb:              cacheDb,
		gitConfig:            appConfig.Global,
		gitOpsClient: gitOps.New(
			logger,
			gitOps.GitConfig{
				GitUser:  appConfig.Global.GitUser,
				GitEmail: appConfig.Global.GitEmail,
			},
			appConfig.GitRepos,
			cacheDb,
		),
	}
	return
}

func loadConfig(logger *zap.SugaredLogger) (appConfig cfg.Config, err error) {
	var rawFile []byte
	if rawFile, err = ioutil.ReadFile(configFile); err == nil {
		if appConfig, err = cfg.ParseConfig(rawFile); err != nil {
			err = errors.Wrap(err, "error parsing config file")
		}
	} else {
		err = errors.Wrap(err, "error reading config")
	}
	if err != nil {
		logger.Errorw("error loading config",
			"laminar.file", configFile,
			"laminar.error", err,
		)
	}
	return
}

func (d *Daemon) Start() {
	go d.webClient.StartWeb()
	go metrics.Start(d.logger)
	d.logger.Debug("opened db: ", configCache)
	d.enterControlLoop()
}

func (d *Daemon) enterControlLoop() {
	ticker := time.Tick(interval)
	for {
		select {
		case repo := <-d.webClient.BuildChan:
			d.singleRepoTask(repo)
		case <-d.webClient.PauseChan:
			d.pause()
		case <-ticker:
			d.masterTask()
		}
	}
}

func (d *Daemon) pause() {
	d.logger.Infow("laminar paused",
		"pauseDuration", pauseDuration,
	)
	<-time.Tick(pauseDuration)
	d.logger.Infow("laminar paused expired. continuing")
}

func (d *Daemon) masterTask() {
	var err error
	err = d.gitOpsClient.UpdateAll()
	if err == nil {
		repoStates := d.gitOpsClient.GetStates()
		for _, state := range repoStates {
			for _, filePath := range state.GetFilesToScan() {
				d.logger.Infow("scanning", filePath)
			}
		}
	}
	if err != nil {
		d.logger.Error(err)
	}
	// from the update policies, make a list of ALL file paths which are referenced in our gitoperations repo
	//for _, state := range d.gitState {
	//	d.updateGitRepoState(state)
	//}

	// TODO: docker reg Timeout?
	// lets gather a full list of docker images we can find matching the configured registries
	//d.scanDockerRegistries()

	// now that we can assume we have some tags in cache, we run a
	// loop over GitRepos
	//for _, state := range d.gitState {
	//	d.updateFiles(*state.repoCfg)
	//}
	if oneShot {
		d.logger.Warn("--one-shot detected.. laminar is now terminating")
		os.Exit(0)
	}
}

func (d *Daemon) singleRepoTask(r web.DockerBuildJSON) {
	//if reg, ok := d.dockerRegistries[r.DockerRegistryUrl]; ok {
	//	for _, state := range d.gitState {
	//		d.updateGitRepoState(state)
	//	}
	//
	//	//d.scanDockerRegistry(reg)
	//
	//	//for _, state := range d.gitState {
	//	//	//d.updateFiles(*state.repoCfg)
	//	//}
	//}
}

//func (d *Daemon) updateFiles(gitRepo cfg.GitRepo) {
//	registryStrings := d.getRegistryStrings()
//	triggerCommitAndPush := false
//	var changes []ChangeRequest
//	for _, updatePolicy := range gitRepo.Updates {
//		var fileList []string
//		// assemble a list of target files for this Update
//		for _, p := range updatePolicy.Files {
//			// get the path of where the gitoperations repo is checked out
//			relativeGitPath := gitoperations.GetRepoPath(gitRepo)
//			// combine these
//			realPath := fmt.Sprintf("%s/%s", relativeGitPath, p.Path)
//
//			// finally this will return all files found
//			for _, paths := range d.opsClient.FindFiles(realPath) {
//				fileList = append(fileList, paths)
//			}
//		}
//
//		for _, filePath := range fileList {
//			d.logger.Debugw("applying update policy",
//				"laminar.file", filePath,
//				"laminar.pattern", updatePolicy.PatternString,
//				"laminar.blacklist", updatePolicy.BlackList,
//			)
//			//newChanges := d.doUpdate(filePath, updatePolicy, registryStrings)
//			if len(newChanges) > 0 {
//				d.logger.Infow("updates desired",
//					"laminar.file", filePath,
//					"laminar.pattern", updatePolicy.PatternString,
//				)
//				triggerCommitAndPush = true
//				for _, stuffDone := range newChanges {
//					changes = append(changes, stuffDone)
//				}
//			}
//		}
//	}
//
//	if triggerCommitAndPush {
//		//d.commitAndPush(changes, gitRepo)
//	}
//}

func (d *Daemon) commitAndPush(changes []shared.ChangeRequest, repo cfg.GitRepo) {
	msg := ""
	if len(changes) > 1 {
		msg = fmt.Sprintf("%s [%d]", d.gitConfig, len(changes))
	} else {
		prettyMessage := gitOps.NicerMessage(changes[0])
		//fmt.Println(changes)
		msg = fmt.Sprintf("%s", prettyMessage)
	}
	d.logger.Infow("doing commit",
		"laminar.gitRepo", repo.URL,
		"laminar.msg", msg,
	)
	d.gitOpsClient.CommitAndPush(repo, msg)
}

//func (d *Daemon) scanDockerRegistries() {
//	for _, dockerReg := range d.dockerRegistries {
//		d.scanDockerRegistry(dockerReg)
//	}
//}

//func (d *Daemon) scanDockerRegistry(dockerReg cfg.DockerRegistry) {
//	d.logger.Infow("scanning docker registry", "url", dockerReg.Reg)
//	foundDockerImages := d.FindDockerImages(
//		d.fileList,
//		fmt.Sprintf(dockerReg.Reg),
//	)
//	if len(foundDockerImages) > 0 {
//		d.dockerRegistryClient.Exec(dockerReg, foundDockerImages)
//		d.logger.Infow("found images (in gitoperations) matching a configured docker registry",
//			"laminar.regName", dockerReg.Name,
//			"laminar.reg", dockerReg.Reg,
//			"laminar.imageCount", len(foundDockerImages),
//		)
//	} else {
//		d.logger.Infow("no images tags found.. ensure the full <image>:<tag> strings present",
//			"laminar.regName", dockerReg.Name,
//			"laminar.reg", dockerReg.Reg,
//		)
//	}
//}

//func (d *Daemon) updateGitRepoState(state gitoperations.RepoState) {
//	// Clone all repos that haven't been cloned yet
//	if state.Repo != nil {
//		d.gitOpsClient.Pull(*state.repoCfg)
//	} else {
//		d.logger.Warnw("repo has not been initialised",
//			"repo.URL", state.repoCfg.URL)
//		//TODO: implement initialisation of uninitialised repos (probs issue for dynamic configs)
//		return
//	}
//
//	// This sections deals with loading remote config from the gitoperations repo
//	// if RemoteConfig is set we want to attempt to read '.laminar.yaml' from the remote repo
//	repoPath := state.repoCfg.GetRealPath()
//	if state.repoCfg.RemoteConfig {
//		d.logger.Debugw("'remote config' == True.. will attempt to update config dynamically",
//			"laminar.repo", state.repoCfg.Name,
//		)
//
//		remoteUpdates, err := cfg.getRemoteUpdates(repoPath, d.logger)
//		if err != nil {
//			d.logger.Warnw("Laminar was told to look at .laminar.yaml but failed",
//				"laminar.repo", state.repoCfg.Name,
//				"laminar.path", repoPath,
//				"laminar.error", err,
//			)
//		}
//
//		// clear out the Updates for this repoNum
//		state.repoCfg.Updates = make([]cfg.Updates, 0)
//		// now assemble that list for this run
//		for _, update := range remoteUpdates.Updates {
//			d.logger.Infow("using 'remote config' from gitoperations repo .laminar.yaml",
//				"laminar.update", update,
//			)
//			state.repoCfg.Updates = append(state.repoCfg.Updates, update)
//		}
//	}
//	// equalise the state.. damn this needs a nice rewrite sometime
//	d.logger.Infow("configured for",
//		"laminar.gitRepo", state.repoCfg.Name,
//		"laminar.updateRules", len(state.repoCfg.Updates),
//	)
//	//d.UpdateFileList(*state.repoCfg)
//
//	// we are ready to dispatch this to start searching the contents of these files
//	d.logger.Debugw("matched files in gitoperations",
//		"laminar.GitRepo", state.repoCfg.Name,
//		"laminar.fileList", d.fileList,
//	)
//}

//func (d *Daemon) getRegistryStrings() []string {
//	// this is a slice of the registry URLs as we expect to see them inside files
//	var registryStrings []string
//	for _, reg := range d.dockerRegistries {
//		registryStrings = append(registryStrings, reg.Reg)
//	}
//	return registryStrings
//}
