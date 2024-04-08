package shutdown

type Runner interface {
	Go(f func() error)
	Wait() error
}
