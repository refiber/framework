package support

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/refiber/framework/constant"
)

func saveTempData(s *support, key constant.SessionKey, data *fiber.Map) error {
	if !key.IsValid() {
		return errors.New("Invalid SessinKey")
	}

	buf, _ := json.Marshal(data)

	session, _ := s.session.Get(s.GetCtx())
	session.Set(string(key)+session.ID(), buf)
	session.SetExpiry(time.Second * 15)

	if err := session.Save(); err != nil {
		return err
	}

	return nil
}

func (s *support) SetSharedData(data fiber.Map) error {
	err := saveTempData(s, constant.SessionKeyShared, &data)
	return err
}

func (s *support) GetSharedData() (*fiber.Map, error) {
	session, err := s.session.Get(s.GetCtx())
	if err != nil {
		return nil, err
	}

	keyShared := string(constant.SessionKeyShared) + session.ID()
	raw := session.Get(keyShared)
	data, ok := raw.([]byte)
	if !ok {
		return nil, nil
	}

	var d fiber.Map
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func GetTempData(session *session.Session) *fiber.Map {
	m := make(fiber.Map)
	m["errors"] = fiber.Map{}
	m["auth"] = new(fiber.Map)
	m["flash"] = new(fiber.Map)
	m["shared"] = new(fiber.Map)

	/**
	 * Form Errors
	 */
	keyErrors := string(constant.SessionKeyError) + session.ID()
	raw := session.Get(keyErrors)
	if data, ok := raw.([]byte); ok {
		var d fiber.Map
		if err := json.Unmarshal(data, &d); err != nil {
			log.Errorw("refiber.support.GetTempData: failed to get errors")
		} else {
			m["errors"] = d
			session.Delete(keyErrors)
		}
	}

	/**
	 * Flash Message
	 */
	keyFlashMessage := string(constant.SessionKeyFlashMessage) + session.ID()
	raw = session.Get(keyFlashMessage)
	if data, ok := raw.([]byte); ok {
		var d fiber.Map
		if err := json.Unmarshal(data, &d); err != nil {
			log.Errorw("refiber.support.GetTempData: failed to get keyFlashMessage")
		} else {
			m["flash"] = d
			session.Delete(keyFlashMessage)
		}
	}

	/**
	 * Auth
	 */
	keyAuth := string(constant.SessionKeyAuth) + session.ID()
	raw = session.Get(keyAuth)
	if data, ok := raw.([]byte); ok {
		var d fiber.Map
		if err := json.Unmarshal(data, &d); err != nil {
			log.Errorw("refiber.support.GetTempData: failed to get auth")
		} else {
			m["auth"] = d
		}
	}

	/**
	 * Shared
	 */
	keyShared := string(constant.SessionKeyShared) + session.ID()
	raw = session.Get(keyShared)
	if data, ok := raw.([]byte); ok {
		var d fiber.Map
		if err := json.Unmarshal(data, &d); err != nil {
			log.Errorw("refiber.support.GetTempData: failed to get shared")
		} else {
			m["shared"] = d
			session.Delete(keyShared)
		}
	}

	if err := session.Save(); err != nil {
		log.Errorw("refiber.support.GetTempData: failed to save session")
	}

	return &m
}
