package hw06pipelineexecution

import (
	"fmt"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outSlise := make([]Out, 0, len(stages))
	outSliseFlag := make([]bool, 0, len(stages))
	result := make(Bi, len(stages)+1)
	defer close(result)
	for task := range in {
		out := make(Bi, 1)
		out <- task
		in = In(out)
		for _, s := range stages {
			in = s(in)
		}
		outSlise = append(outSlise, in)
		outSliseFlag = append(outSliseFlag, false)
	}
	for {
		select {
		case _, ok := <-done:
			if !ok {
				return result
			}
		default:
		}
		for i := 0; i < len(outSlise); i++ {
			select {
			case x := <-outSlise[i]:
				outSliseFlag[i] = true
				result <- x
				fmt.Println(x)
			default:
			}
		}
		quit := true
		for _, f := range outSliseFlag {
			quit = quit && f
		}
		if quit {
			break
		}
	}
	return result
}
