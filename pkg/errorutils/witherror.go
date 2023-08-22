package errorutils

func Drop[T interface{}](res T, err error) T {
	return res
}

func Panic[T interface{}](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res
}

func WithError[T interface{}](fn func(err error)) func(T, error) T {
	return func(res T, err error) T {
		if err != nil {
			fn(err)
		}
		return res
	}
}
