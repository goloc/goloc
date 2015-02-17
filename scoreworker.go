package goloc

import ()

func scoreWorker(result *Result, scorer Scorer, jobChan <-chan bool, resultChan chan<- *Result) {
	result.Score = scorer(result)
	resultChan <- result
	<-jobChan
}
