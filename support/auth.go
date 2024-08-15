package support

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/csrf"

	"github.com/refiber/framework/constant"
)

func (s *support) Auth(ctx *fiber.Ctx) *auth {
	return &auth{s, ctx}
}

type auth struct {
	support *support
	ctx     *fiber.Ctx
}

func (a *auth) NewAuthenticatedUserSession(user interface{}) error {
	session, err := a.support.sessionStore.Get(a.ctx)
	if err != nil {
		return err
	}

	// before destroy the session, get redirection & crsf token from the session, then save it to the new session
	var redirection *string
	if v := session.Get(string(constant.SessionKeyRedirection)); v != nil {
		_v := v.(string)
		redirection = &_v
	}

	var crsfToken *csrf.Token
	if v := session.Get(string(constant.SessionKeyCSRFToken)); v != nil {
		_v := v.(csrf.Token)
		crsfToken = &_v
	}

	session.Destroy()

	// get new session
	sessionNew, err := a.support.sessionStore.Get(a.ctx)
	if err != nil {
		return err
	}

	// TODO: make this customizable
	sessionNew.SetExpiry((time.Hour * 24) * 7)

	if redirection != nil && *redirection != "" {
		sessionNew.Set(string(constant.SessionKeyRedirection), redirection)
	}

	if crsfToken != nil {
		sessionNew.Set(string(constant.SessionKeyCSRFToken), *crsfToken)
	}

	sessionKey := string(constant.SessionKeyAuth)
	buf, _ := json.Marshal(user)
	sessionNew.Set(sessionKey, buf)

	if err := sessionNew.Save(); err != nil {
		return err
	}

	return nil
}

func (a *auth) GetAuthenticatedUserSession(user interface{}) error {
	session, err := a.support.sessionStore.Get(a.ctx)
	if err != nil {
		return err
	}

	raw := session.Get(string(constant.SessionKeyAuth))

	if data, ok := raw.([]byte); ok {
		if err := json.Unmarshal(data, &user); err != nil {
			return err
		}
	}

	return nil
}

func (a *auth) UpdateAuthenticatedUserSession(user interface{}) error {
	session, err := a.support.sessionStore.Get(a.ctx)
	if err != nil {
		return err
	}

	buf, _ := json.Marshal(user)
	session.Set(string(constant.SessionKeyAuth), buf)

	if err := session.Save(); err != nil {
		return err
	}

	return nil
}

func (a *auth) DestroyAuthenticatedUserSession() error {
	session, err := a.support.sessionStore.Get(a.ctx)
	if err != nil {
		return err
	}

	session.Destroy()

	return nil
}

// get protected url
func getRedirectLocation(a *auth) (location *string, err error) {
	if session, err := a.support.sessionStore.Get(a.ctx); err == nil {
		key := string(constant.SessionKeyRedirection)
		data := session.Get(key)
		if redirectLocation, ok := data.(string); ok && redirectLocation != "" {
			session.Delete(key)
			if err := session.Save(); err != nil {
				return nil, err
			}
			return &redirectLocation, nil
		}
	}

	return nil, err
}

func (a *auth) RedirectTo(defaultLocation string) error {
	var l *string
	l = &defaultLocation

	if redirectLocation, err := getRedirectLocation(a); redirectLocation != nil {
		if err != nil {
			log.Errorf("refiber.support.auth.RedirectTo:", err)
		}

		l = redirectLocation
	}

	return a.support.Redirect(a.ctx).To(*l).Now()
}

func (a *auth) RedirectToWithMessage(defaultLocation string, messageType MessageType, message string) error {
	var l *string
	l = &defaultLocation

	if redirectLocation, err := getRedirectLocation(a); redirectLocation != nil {
		if err != nil {
			log.Errorf("refiber.support.auth.RedirectToWithMessage:", err)
		}

		l = redirectLocation
	}

	return a.support.Redirect(a.ctx).To(*l).WithMessage(messageType, message).Now()
}

func (a *auth) LoginPage(location string) error {
	// save protected url, then redirect back to the protected url after login
	if session, err := a.support.GetSessionStore().Get(a.ctx); err == nil {
		session.Set(string(constant.SessionKeyRedirection), a.ctx.OriginalURL())
		if err := session.Save(); err != nil {
			log.Errorw("refiber.support.auth.LoginPage: failed to save session")
		}
	}

	return a.ctx.Redirect(location)
}
