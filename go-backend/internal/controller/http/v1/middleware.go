package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	cuserr "github.com/kurochkinivan/pulskrsk/internal/customErrors"

	"github.com/sirupsen/logrus"
)

const (
	RequestIDKey = 0
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func logMdw(next appHandler) appHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		entry := logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
			"request_id":  middleware.GetReqID(r.Context()),
		})

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t := time.Now()

		defer func() {
			entry.WithFields(logrus.Fields{
				"status":   ww.Status(),
				"size":     ww.BytesWritten(),
				"duration": time.Since(t),
			}).Info("request completed")
		}()

		return next(ww, r)
	}
}

func authMdw(next appHandler, signingKey string) appHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		accessTokenCookie, err := r.Cookie("access_token")
		if err != nil {
			return cuserr.ErrGetCookie.WithErr(err)
		}
		accessToken := accessTokenCookie.Value

		payload, err := parseToken(signingKey, accessToken)
		if err != nil {
			return err
		}

		exp := time.Unix(int64(payload["exp"].(float64)), 0)
		if time.Now().After(exp) {
			return cuserr.ErrTokenExired
		}

		r.Header.Set("user_id", payload["ueid"].(string))

		return next(w, r)
	}
}

func errMdw(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customerr *cuserr.AppError
		err := h(w, r)

		if err != nil {
			if errors.As(err, &customerr) {
				w.WriteHeader(customerr.HTTPCode)
				w.Write(customerr.MarshalWithTrace(err.Error()))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(cuserr.SystemError(err, "appHandler.errMdw", "not custom error").MarshalWithTrace(err.Error()))
		}
	}
}

func parseToken(signingKey, accessToken string) (jwt.MapClaims, error) {
	logrus.WithField("token", accessToken).Debug("parsing jwt-token")
	const op string = "middleware.ParseToken"

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.WithField("alg", t.Header["alg"]).Error("unexpected signing method")
			return nil, cuserr.SystemError(nil, op, "unexpected signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, cuserr.ErrTokenExired
		}
		return nil, cuserr.ErrInvalidAuthHeader.WithErr(fmt.Errorf("%s: %w", op, err))
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, cuserr.SystemError(nil, op, "claims are not of type jwt.MapClaims")
	}

	return payload, nil
}
