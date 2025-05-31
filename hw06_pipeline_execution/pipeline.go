package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, myStage := range stages {
		out = myStage(checkStage(out, done))
	}
	return out
}

func checkStage(in In, done In) Out {
	res := make(Bi)
	go func() {
		defer close(res)
		for {
			select {
			case <-done:
				return
			case inVal, isOk := <-in:
				if !isOk {
					return
				}
				res <- inVal
			}
		}
	}()
	return res
}
