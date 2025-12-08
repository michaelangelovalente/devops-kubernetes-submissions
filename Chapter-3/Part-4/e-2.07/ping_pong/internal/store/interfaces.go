package store

type PingPongRepo interface {
	Update() (int, error)
	GetCurr() (int, error)
}
