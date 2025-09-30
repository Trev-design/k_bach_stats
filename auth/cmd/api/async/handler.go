package async

type Task chan Output

func NewTask(input Input, fn Fn) Task {
	task := make(chan Output)
	go fn(task, input)
	return task
}

func (task Task) AwaitResult() Output {
	return <-task
}
