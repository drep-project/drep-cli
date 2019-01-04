package util

type MyError struct {
    Err error
}

func (e *MyError) Error() string {
    return e.Err.Error()
}

type DataError struct {
    MyError
}

type TimeoutError struct {
    MyError
}

type ConnectionError struct {
    MyError
}

type TransmissionError struct {
    MyError
}

type DefaultError struct {
    MyError
}

type DupOpError struct {
    MyError
}

type OfflineError struct {
    MyError
}

type ConsensusError struct {
    MyError
}