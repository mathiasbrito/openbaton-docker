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

// sessionEntry is uniquely associated with a set of Credentials, and
// holds a section for it. Adding a mutex here avoids races while reconnecting.
type sessionEntry struct {
	sess *session
	mux sync.Mutex
}

type sessionCache struct {
	// getEntry() is a critical section
	lock sync.Mutex

	// Use creds.Credentials (a host,user,password tuple) as a key for a session.
	// A session authenticated for the given credentials on the given host should be valid.
	sessions map[creds.Credentials]*sessionEntry
}

// get() retrieves or creates a session for the given credentials.
// It requires a mutex to avoid multiple parallel get() requests.
func (sc *sessionCache) get(c creds.Credentials) (*session, error) {
	entry := sc.getEntry(c)

	entry.mux.Lock()
	defer entry.mux.Unlock()

	sess := entry.sess
	if sess != nil && !sess.invalid {
		return sess, nil
	}

	sess, err := newSession(c)
	if err != nil {
		return nil, err
	}

	// an invalid session object will be overwritten here.
	entry.sess = sess

	return sess, nil
}

func (sc *sessionCache) getEntry(c creds.Credentials) *sessionEntry {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	entry, found := sc.sessions[c]
	if !found {
		entry = &sessionEntry{}
		sc.sessions[c] = entry
	}

	return entry
}

func (sc *sessionCache) init() {
	sc.sessions = make(map[creds.Credentials]*sessionEntry)
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

	cache.lock.Lock()
	defer cache.lock.Unlock()

	for _, sessEntry := range cache.sessions {

		if sess := sessEntry.sess; sess != nil {
			if err := sess.logout(); err != nil {
				ret = append(ret, err)
			}

			if err := sess.conn.Close(); err != nil {
				ret = append(ret, err)
			}
		}
	}

	cache.init()

	if len(ret) == 0 {
		return nil
	}

	return ret
}
