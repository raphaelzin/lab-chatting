package accessors

import (
	"main/models"
)

type UserAccessor struct {
	BaseAccessor[models.User]
}

func (a UserAccessor) Save(user models.User) error {
	// Save on DB, if it's successful, save on cache
	return a.cacheSave(user.Id, user)
}

func (a UserAccessor) GetById(id string) (models.User, error) {
	return a.cacheGet(id)
}
