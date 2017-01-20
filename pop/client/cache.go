package client

import (
	"bytes"
	"fmt"
	"sync"
)

var cache sessionCache

func init() {
	cache.init()
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

func (sc *sessionCache) init() {
	sc.sessions = make(map[Credentials]*session)
}

// FlushError is an error type containing any error encountered while executing 
// FlushSessions.
type FlushError []error

func (errs FlushError) Error() string {
	buf := bytes.NewBufferString("got errors while closing sessions: ")

	for _, err := range errs {
		buf.WriteString(fmt.Sprintf(`"%v" `, err))
	}

	return buf.String()
}

// FlushSessions closes all the cached sessions, trying to log out from them first.
// If you wish, you can inspect the single errors from the Logout and Close operation
// through type casting the returned error into a FlushError.
func FlushSessions() error {
	ret := FlushError{}

	for _, sess := range cache.sessions {
		if err := sess.logout(); err != nil {
			ret = append(ret, err)
		}

		if err := sess.conn.Close(); err != nil {
			ret = append(ret, err)
		}
	}

	cache.init()

	return ret
}
