package constant

type SessionKey string

const (
	SessionKeyError        SessionKey = "form:errors:"
	SessionKeyAuth         SessionKey = "auth:"
	SessionKeyRedirection  SessionKey = "redirection:"
	SessionKeyFlashMessage SessionKey = "flash_message:"
	SessionKeyShared       SessionKey = "shared:"
)

func (k SessionKey) IsValid() bool {
	switch k {
	case SessionKeyError, SessionKeyAuth, SessionKeyRedirection, SessionKeyFlashMessage, SessionKeyShared:
		return true
	}

	return false
}
