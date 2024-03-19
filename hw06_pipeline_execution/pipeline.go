package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outSlise := make([]Out, 0)
	outSliseFlag := make([]bool, 0)
	resultSlice := make([]interface{}, 0)
	for task := range in {
		out := make(Bi, 1)
		out <- task
		in = In(out)
		for _, s := range stages {
			in = s(in)
		}
		outSlise = append(outSlise, in)
		outSliseFlag = append(outSliseFlag, false)
		resultSlice = append(resultSlice, nil)
	}
	result := make(Bi, len(resultSlice))
	defer close(result)
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
				resultSlice[i] = x
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
	for _, v := range resultSlice {
		result <- v
	}
	return result
}
