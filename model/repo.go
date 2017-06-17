package model

import "log"

type UserRepository interface {
	FindByID(id uint) (User, error)
}

type DummyUserRepository struct{}

func (r *DummyUserRepository) FindByID(id uint) (User, error) {
	return User{
		ID:     id,
		Handle: "user",
		Email:  "email@test.com",
	}, nil
}

type CollectionRepository interface {
	FindByID(id uint) (Collection, error)
}

type DummyCollectionRepository struct{}

func (r *DummyCollectionRepository) FindByID(id uint) (Collection, error) {

	col := Collection{
		Record: Record{
			ID: id,
		},
	}
	col.UpdateSlug("bam")
	return col, nil
}

func (r *DummyCollectionRepository) Save(col Collection) (Collection, error) {
	log.Printf("Saving collection %s", col)
	return col, nil
}
