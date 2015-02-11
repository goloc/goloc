package goloc

import ()

func scoreWorker(result *Result, jobChan <-chan bool, resultChan chan<- *Result) {
	s := Score(result.Search, result.Name) - int(result.Priority)
	result.Score = s
	resultChan <- result
	<-jobChan
}
