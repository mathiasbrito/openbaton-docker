package client

import "sync"

var cache sessionCache

func init() {
	cache.sessions = make(map[Credentials]*session)
}

type sessionCache struct {
	lock sync.Mutex

	sessions map[Credentials]*session
}

func (sc *sessionCache) get(c Credentials) (*session, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	sess, ok := sc.sessions[c]
	if ok && !sess.invalid {
		return sess, nil
	}

	sess, err := newSession(c)
	if err != nil {
		return nil, err
	}

	// an invalid session object will be overwritten here.
	sc.sessions[c] = sess

	return sess, nil
}
