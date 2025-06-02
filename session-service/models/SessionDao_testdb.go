package models

type TestDBRepo struct {
}

func (p *TestDBRepo) CreateSession(sessionDao *Session) (*Session, error) {
	return nil, nil
}

func (p *TestDBRepo) StopSession(sessionDao *Session) (*Session, error) {
	return nil, nil
}
