package client

import (
	"bytes"
	"fmt"
	"sync"
	
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
)

// global session cache.
var cache sessionCache

// The cache contains a map, so it must be initialised at startup.
func init() {
	cache.init()
}

type sessionCache struct {
	// get() is a critical section
	lock sync.Mutex

	// Use creds.Credentials (a host,user,password tuple) as a key for a session.
	// A session authenticated for the given credentials on the given host should be valid.
	sessions map[creds.Credentials]*session
}

// get() retrieves or creates a session for the given credentials.
// It requires a mutex to avoid multiple parallel get() requests. 
func (sc *sessionCache) get(c creds.Credentials) (*session, error) {
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
	sc.sessions = make(map[creds.Credentials]*session)
}

// FlushError is an error type containing any error encountered while executing
// FlushSessions.
type FlushError []error

// Error returns a string that summarises the various errors contained by this FlushError.
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

	if len(ret) == 0 {
		return nil
	}

	return ret
}
