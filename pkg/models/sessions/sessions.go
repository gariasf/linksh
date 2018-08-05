package sessions

import (
    "errors"
	"time"
)

//Session is the datatype which contains all the relevant information about a session
type Session struct {
	ID        string
	OwnerID   string
	ExpiresOn int64
	CreatedAt int64
}

//SessionManager is helps to manage the session stored in the provider
type SessionManager struct {
    provider SessionProvider
    autoGCShouldBeRunning bool
}

//NewSessionManager Returns a new SessionManager with the selected provider
func NewSessionManager(provider SessionProvider, providerArguments ...interface{}) *SessionManager {
	manager := &SessionManager{provider: provider}
	return manager
}

//Add a session to the provider's storage
func (manager *SessionManager) Add(session Session) error {
   return manager.provider.Add(session)
}

//Get returns the requested session from the provider's storage, if not found it returns an error
func (manager *SessionManager) Get(id string) (Session, error) {
    return manager.provider.Get(id)
}

//GetByOwnerID Returns the sessions which belongs to the selected user
func (manager *SessionManager) GetByOwnerID(ownerID string) (map[string]Session, error) {
    return manager.provider.GetByOwnerID(ownerID)
}

//Update a session
func (manager *SessionManager) Update(session Session) error {
    return manager.provider.Update(session)
}

//Delete the session with selected id
func (manager *SessionManager) Delete(id string) error {
    return manager.provider.Delete(id)
}

//GC deletes expired entries from the provider's storage
func (manager *SessionManager) GC() error {
    return manager.provider.GC()
}

//EnableAutoGC start a background job which will run the GC function every time "x" has passed
// Only a job of autoGC can be running at the same time for each SessionManager instance.
func (manager *SessionManager) EnableAutoGC(x time.Duration) error {
    if manager.autoGCShouldBeRunning {
        return errors.New("The autoGC job is already running")
    }
    manager.autoGCShouldBeRunning = true
    go func() {
        for manager.autoGCShouldBeRunning {
            manager.GC()
            time.Sleep(x)
        }
    }()

    return nil
}

//DisableAutoGC stops the autoGC job.
func (manager *SessionManager) DisableAutoGC() error {
    if !manager.autoGCShouldBeRunning {
        return errors.New("The autoGC job was not running")
    }
    manager.autoGCShouldBeRunning = false

    return nil
}


//SessionProvider is the interface for a valid session storage
type SessionProvider interface {
	Add(session Session) error
	Get(id string) (Session, error)
	GetByOwnerID(ownerID string) (map[string]Session, error)
	Update(session Session) error
	Delete(id string) error
	GC() error
}
