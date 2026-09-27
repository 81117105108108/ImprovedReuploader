package ide

import (
	"bytes"
	"errors"

	"github.com/kartFr/Asset-Reuploader/internal/roblox"
)

var UploadAnimationErrors = struct {
	ErrNotLoggedIn       error
	ErrTokenInvalid      error
	ErrInappropriateName error
}{
	ErrNotLoggedIn:       errors.New("not logged in"),
	ErrTokenInvalid:      errors.New("XSRF token validation failed"),
	ErrInappropriateName: errors.New("inappropriate name or description"),
}

func NewUploadAnimationHandler(
	c *roblox.Client,
	name,
	description string,
	data *bytes.Buffer,
	groupID ...int64,
) (func() (int64, error), error) {
	return NewUploadHandler(c, "Animation", "model/x-rbxm", name, description, data, UploadAnimationErrors.ErrTokenInvalid, UploadAnimationErrors.ErrNotLoggedIn, groupID...)
}
