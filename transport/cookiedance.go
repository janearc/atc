package transport

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// cookie dance is a set of methods to abstract the interaction between backends
// like strava and the service. a single call will return a sufficiently authenticated
// object to make subsequent requests to the backend.

// FLOW:
// 1. we handle a request
// 2. we check for a cookie
// 3. if no cookie, redirect to auth
// 4. if cookie, check if expired
// 5. if expired, refresh token
// 6. if refreshed, update cookie

// since these requests happen in the backend, we shouldn't have to pass tokens to functions
// in the backend. so this might look like:

func (t *Transport) CookieUp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("[%s]: Received request to update credentials", r.URL.Path)

		// check to see if we have a code in the request
		code := r.URL.Query().Get("code")
		if code == "" {
			logrus.Error("Code not found in callback")
			http.Error(w, "Code not found", http.StatusBadRequest)

			// okay this needs to head out to strava and get oauth done, we can't really help you here
			return
		} else {
			err := t.ExchangeCodeForToken(code)
			if err != nil {
				// we had a code, but for some reason we got an error back. this is unrecoverable.
				logrus.WithError(err).Errorf("[%s]: Invalid code provided [%s], auth unrecoverable.", r.URL.Path, code)
			} else {
				// we send strava a code, which should have returned to us a token which we
				// can use for additional requests.
				logrus.Infof("[%s]: successfully exchanged code [%s] for new authentication token", r.URL.Path, code)
				if t.GetAccessToken() != "" {
					// we have a token, let's update the cookie
					logrus.Infof("[%s]: Successfully retrieved access token, writing cookies", r.URL.Path)

					// Cookie the access token
					http.SetCookie(w, &http.Cookie{
						Name:       "strava_token",
						Value:      t.GetAccessToken(),
						Path:       "/",
						Domain:     "",
						Expires:    time.Now().Add(24 * time.Hour),
						RawExpires: "",
						MaxAge:     0,
						Secure:     true,
						HttpOnly:   true,
						SameSite:   0,
						Raw:        "",
						Unparsed:   nil,
					})

					// Cookie the refresh token
					http.SetCookie(w, &http.Cookie{
						Name:     "strava_refresh_token",
						Value:    t.GetRefreshToken(),
						Path:     "/",
						Expires:  time.Now().Add(30 * 24 * time.Hour), // Refresh token typically lasts longer
						HttpOnly: true,
						Secure:   true,
					})

					// Redirect to the home page after setting the cookies
					logrus.Infof("[%s]: Redirecting to / with cookies", r.URL.Path)

					// update the backend struct to let the service know we tried
					t.authenticated = true

					// send them back to the app
					http.Redirect(w, r, "/", http.StatusFound)
				}
			}
		}
	})

	logrus.Info("Returning from CookieUp")
	return
}
