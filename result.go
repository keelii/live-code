package main

type Result[T any] struct {
	Ok   bool   `json:"ok"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func OkResult[T any](data T) Result[T] {
	return Result[T]{
		Ok:   true,
		Msg:  "",
		Data: data,
	}
}
func ErrResult(msg string) Result[string] {
	return Result[string]{
		Ok:   false,
		Msg:  msg,
		Data: "",
	}
}
