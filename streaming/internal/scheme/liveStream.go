package scheme

type LiveStream struct {
	ID     string
	Title  string
	IsLive int

	Key string
}

type Filter struct {
	Key    string
	Id     string
	Title  string
	IsLive bool
}

type LiveStreamsErr struct {
	cErr <-chan error
	c    <-chan *LiveStream
	err  error
}

func LiveStreamIterator(c <-chan *LiveStream, cErr <-chan error) *LiveStreamsErr {
	return &LiveStreamsErr{
		cErr: cErr,
		c:    c,
		err:  nil,
	}
}

func (l *LiveStreamsErr) Next() (*LiveStream, bool) {
	ls, ok := <-l.c
	return ls, ok
}

func (l *LiveStreamsErr) Err() error {
	return <-l.cErr
}
