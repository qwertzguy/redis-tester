package main

import (
	"fmt"

	"math/rand"
	"time"
)

// StageRunnerResult is returned from StageRunner.Run()
type StageRunnerResult struct {
	lastStageIndex int
	error          error
	logger         *customLogger
}

// IsSuccess says whether a StageRunnerResult was successful
// or not
func (res StageRunnerResult) IsSuccess() bool {
	return res.error == nil
}

// StageRunner is used to run multiple stages
type StageRunner struct {
	stages  []Stage
	isDebug bool
}

func newStageRunner(isDebug bool) StageRunner {
	return StageRunner{
		isDebug: isDebug,
		stages: []Stage{
			Stage{
				name:    "Stage 1: Bind to a port",
				logger:  getLogger(isDebug, "[stage-1] "),
				runFunc: testBindToPort,
			},
			Stage{
				name:    "Stage 2: PING <-> PONG",
				logger:  getLogger(isDebug, "[stage-2] "),
				runFunc: testPingPong,
			},
			Stage{
				name:    "Stage 3: ECHO... O... O...",
				logger:  getLogger(isDebug, "[stage-3] "),
				runFunc: testEcho,
			},
			Stage{
				name:    "Stage 4: Multiple Clients",
				logger:  getLogger(isDebug, "[stage-4] "),
				runFunc: testMultipleClients,
			},
			Stage{
				name:    "Stage 5: SET & GET",
				logger:  getLogger(isDebug, "[stage-5] "),
				runFunc: testGetSet,
			},
			Stage{
				name:    "Stage 6: Expiry!",
				logger:  getLogger(isDebug, "[stage-6] "),
				runFunc: testExpiry,
			},
		},
	}
}

// Run tests in a specific StageRunner
func (r StageRunner) Run(executable *Executable) StageRunnerResult {
	for index, stage := range r.stages {
		logger := stage.logger
		logger.Infof("Running test: %s", stage.name)

		stageResultChannel := make(chan error, 1)
		go func() {
			err := stage.runFunc(executable, logger)
			stageResultChannel <- err
		}()

		var err error
		select {
		case stageErr := <-stageResultChannel:
			err = stageErr
		case <-time.After(5 * time.Second):
			err = fmt.Errorf("timed out, test exceeded 5 seconds")
		}

		if err != nil {
			reportTestError(err, r.isDebug, logger)
			return StageRunnerResult{
				lastStageIndex: index,
				error:          err,
			}
		}

		logger.Successf("Test passed.")
	}

	return StageRunnerResult{
		lastStageIndex: len(r.stages) - 1,
		error:          nil,
	}
}

// Truncated returns a stageRunner with fewer stages
func (r StageRunner) Truncated(stageIndex int) StageRunner {
	maxStageIndex := min(stageIndex, len(r.stages)-1)
	newStages := make([]Stage, maxStageIndex+1)
	for i := 0; i <= maxStageIndex; i++ {
		newStages[i] = r.stages[i]
	}

	return StageRunner{
		isDebug: r.isDebug,
		stages:  newStages,
	}
}

// Fuck you, go
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Randomized returns a stage runner that has stages randomized
func (r StageRunner) Randomized() StageRunner {
	return StageRunner{
		isDebug: r.isDebug,
		stages:  shuffleStages(r.stages),
	}
}

func shuffleStages(stages []Stage) []Stage {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]Stage, len(stages))
	perm := r.Perm(len(stages))
	for i, randIndex := range perm {
		ret[i] = stages[randIndex]
	}
	return ret
}

func reportTestError(err error, isDebug bool, logger *customLogger) {
	logger.Errorf("%s", err)
	if isDebug {
		logger.Errorf("Test failed")
	} else {
		logger.Errorf("Test failed " +
			"(try setting 'debug: true' in your codecrafters.yml to see more details)")
	}
}

// Stage is blah
type Stage struct {
	name        string
	description string
	runFunc     func(executable *Executable, logger *customLogger) error
	logger      *customLogger
}
