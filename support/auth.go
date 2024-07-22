package support

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/refiber/framework/constant"
)

func (s *support) NewAuthenticatedUserSession(user interface{}) error {
	session, err := s.session.Get(s.GetCtx())
	if err != nil {
		return err
	}

	// before destroy the session, get redirection from the session, then save redirection to the new session
	var redirection string
	_redirection := session.Get(string(constant.SessionKeyRedirection) + session.ID())
	if _redirection != nil {
		redirection = _redirection.(string)
	}
	session.Destroy()

	// get new session
	sessionNew, err := s.session.Get(s.GetCtx())
	if err != nil {
		return err
	}

	sessionNew.SetExpiry((time.Hour * 24) * 7)

	if redirection != "" {
		sessionNew.Set(string(constant.SessionKeyRedirection)+sessionNew.ID(), redirection)
	}

	sessionKey := string(constant.SessionKeyAuth) + sessionNew.ID()
	buf, _ := json.Marshal(user)
	sessionNew.Set(sessionKey, buf)

	if err := sessionNew.Save(); err != nil {
		return err
	}

	return nil
}

func (s *support) GetAuthenticatedUserSession(user interface{}) error {
	session, err := s.session.Get(s.GetCtx())
	if err != nil {
		return err
	}

	raw := session.Get(string(constant.SessionKeyAuth) + session.ID())

	if data, ok := raw.([]byte); ok {
		if err := json.Unmarshal(data, &user); err != nil {
			return err
		}
	}

	return nil
}

func (s *support) UpdateAuthenticatedUserSession(user interface{}) error {
	session, err := s.session.Get(s.GetCtx())
	if err != nil {
		return err
	}

	buf, _ := json.Marshal(user)
	session.Set(string(constant.SessionKeyAuth)+session.ID(), buf)

	if err := session.Save(); err != nil {
		return err
	}

	return nil
}

func (s *support) DestroyAuthenticatedUserSession() error {
	session, err := s.session.Get(s.GetCtx())
	if err != nil {
		return err
	}

	session.Reset()

	return nil
}

func getRedirectLocation(s Refiber) (location *string, err error) {
	if session, err := s.GetSession().Get(s.GetCtx()); err == nil {
		key := string(constant.SessionKeyRedirection) + session.ID()
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

func AuthRedirection(s Refiber, location string) error {
	var l *string
	l = &location

	if redirectLocation, err := getRedirectLocation(s); redirectLocation != nil {
		if err != nil {
			log.Errorf("refiber.auth.AuthRedirection:", err)
		}

		l = redirectLocation
	}

	return s.Redirect().To(*l).Now()
}

func AuthRedirectionWithMessage(s Refiber, location string, messageType MessageType, message string) error {
	var l *string
	l = &location

	if redirectLocation, err := getRedirectLocation(s); redirectLocation != nil {
		if err != nil {
			log.Errorf("refiber.auth.AuthRedirection:", err)
		}

		l = redirectLocation
	}

	return s.Redirect().To(*l).WithMessage(messageType, message).Now()
}

func AuthLoginPage(location string, s Refiber) error {
	if session, err := s.GetSession().Get(s.GetCtx()); err == nil {
		session.Set(string(constant.SessionKeyRedirection)+session.ID(), s.GetCtx().OriginalURL())
		if err := session.Save(); err != nil {
			log.Errorw("refiber.support.auth.AuthLoginPage: failed to save session")
		}
	}

	return s.GetCtx().Redirect(location)
}
