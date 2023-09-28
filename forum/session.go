package forum

// // Session represents a user session
// type Session struct {
// 	UserID      string
// 	SessionID   string // Add the SessionID field
// 	Expiration  time.Time
// 	LastUpdated time.Time
// }

// // sessions maps user IDs to Session objects
// var sessions = make(map[string]*Session)
// var mu sync.Mutex // Mutex for concurrent access to the sessions map

// // sessionDuration is the duration of a session (e.g., 30 minutes)
// const sessionDuration = 30 * time.Minute

// // sessionCleanupInterval is the interval at which expired sessions are cleaned up
// const sessionCleanupInterval = 1 * time.Hour

// // init initializes the session cleanup process in the background
// func init() {
// 	go cleanupExpiredSessions()
// }

// // createSession creates a new session for the given user ID
// func createSession(userID string) (*Session, error) {
// 	// Lock the mutex for safe concurrent access
// 	mu.Lock()
// 	defer mu.Unlock()

// 	// Check if a session already exists for the user
// 	existingSession := sessions[userID]
// 	if existingSession != nil {
// 		// Invalidate the existing session
// 		invalidateSession(existingSession)
// 	}

// 	// Set the session expiration time
// 	expiration := time.Now().Add(sessionDuration)

// 	// Create a new session
// 	session := &Session{
// 		UserID:      userID,
// 		Expiration:  expiration,
// 		LastUpdated: time.Now(),
// 	}

// 	// Store the session in the sessions map
// 	sessions[userID] = session

// 	return session, nil

// }

// // getSession retrieves a session by user ID
// func getSession(userID string) *Session {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	return sessions[userID]
// }

// // invalidateSession invalidates a session by removing it from the sessions map
// func invalidateSession(session *Session) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	delete(sessions, session.UserID)
// }

// // cleanupExpiredSessions periodically removes expired sessions
// func cleanupExpiredSessions() {
// 	for {
// 		time.Sleep(sessionCleanupInterval)

// 		// Current time
// 		now := time.Now()

// 		mu.Lock()
// 		// Iterate over sessions and remove expired ones
// 		for userID, session := range sessions {
// 			if now.After(session.Expiration) {
// 				// Session has expired, remove it
// 				delete(sessions, userID)
// 			}
// 		}
// 		mu.Unlock()
// 	}
// }

// // setSessionCookie sets the session ID in the response cookie
// func setSessionCookie(w http.ResponseWriter, userID string) {
// 	cookie := http.Cookie{
// 		Name:     "session",
// 		Value:    userID,
// 		Expires:  time.Now().Add(sessionDuration),
// 		HttpOnly: true,
// 	}
// 	http.SetCookie(w, &cookie)
// }

// // extractSessionIDFromCookie extracts the session ID from the request cookie
// func extractSessionIDFromCookie(r *http.Request) (string, error) {
// 	cookie, err := r.Cookie("session")
// 	if err != nil {
// 		return "", err
// 	}
// 	return cookie.Value, nil
// }
