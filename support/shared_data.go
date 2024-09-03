package support

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/refiber/framework/constant"
)

func (s *support) SharedData(ctx *fiber.Ctx) *sharedData {
	return &sharedData{s, ctx}
}

type sharedData struct {
	support *support
	ctx     *fiber.Ctx
}

func (s *sharedData) saveTempData(key constant.SessionKey, data *fiber.Map) error {
	if !key.IsValid() {
		return errors.New("Invalid SessinKey")
	}

	buf, _ := json.Marshal(data)

	session, _ := s.support.sessionStore.Get(s.ctx)
	session.Set(string(key), buf)
	session.SetExpiry(time.Second * 15)

	if err := session.Save(); err != nil {
		return err
	}

	return nil
}

func (s *sharedData) Set(data fiber.Map) error {
	err := s.saveTempData(constant.SessionKeyShared, &data)
	return err
}

func (s *sharedData) Get() (*fiber.Map, error) {
	session, err := s.support.sessionStore.Get(s.ctx)
	if err != nil {
		return nil, err
	}

	keyShared := string(constant.SessionKeyShared)
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

func (s *sharedData) GetTemp() *fiber.Map {
	m := make(fiber.Map)
	m["errors"] = fiber.Map{}
	m["auth"] = new(fiber.Map)
	m["flash"] = new(fiber.Map)
	m["shared"] = new(fiber.Map)

	session, err := s.support.sessionStore.Get(s.ctx)
	if err != nil {
		log.Errorf("refiber.support.session.GetTempData:", err)
		return &m
	}

	/**
	 * Form Errors
	 */
	keyErrors := string(constant.SessionKeyError)
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
	keyFlashMessage := string(constant.SessionKeyFlashMessage)
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
	keyAuth := string(constant.SessionKeyAuth)
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
	keyShared := string(constant.SessionKeyShared)
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
