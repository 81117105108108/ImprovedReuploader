package ide

import (
	"bytes"
	"errors"

	"github.com/kartFr/Asset-Reuploader/internal/roblox"
)

var UploadMeshErrors = struct {
	ErrNotLoggedIn       error
	ErrTokenInvalid      error
	ErrInappropriateName error
}{
	ErrNotLoggedIn:       errors.New("not logged in"),
	ErrTokenInvalid:      errors.New("XSRF token validation failed"),
	ErrInappropriateName: errors.New("inappropriate name or description"),
}

func NewUploadMeshHandler(
	c *roblox.Client,
	name,
	description string,
	data *bytes.Buffer,
	groupID ...int64,
) (func() (int64, error), error) {
	return NewUploadHandler(c, "Mesh", "model/x-file-mesh-data", name, description, data, UploadMeshErrors.ErrTokenInvalid, UploadMeshErrors.ErrNotLoggedIn, groupID...)
}
