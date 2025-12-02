package async

type result[O any] struct {
	value O
	err   error
}

type Task[O any] chan result[O]

func NewTask[I any, O any](input I, fn func(I) (O, error)) Task[O] {
	task := make(chan result[O])

	go func() {
		output, err := fn(input)
		task <- result[O]{
			value: output,
			err:   err,
		}
	}()

	return task
}

func (task Task[O]) AwaitResult() (O, error) {
	result := <-task
	return result.value, result.err
}
