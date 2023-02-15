package test

import (
	"errors"
	"lsat/auth"
	. "lsat/secrets"
)

type UserProfile struct {
	Id     UserId
	Root   Secret
	tokens []auth.LSAT
}

func NewUserProfile(Id UserId, Root Secret) UserProfile {
	return UserProfile{
		Id,
		Root,
		make([]auth.LSAT, 10),
	}
}

type UserBase struct {
	minter auth.Minter
	users  map[UserId]UserProfile
}

func (users UserBase) GetRoot(uid UserId) (Secret, error) {
	profile, ok := users.users[uid]
	if ok {
		return profile.Root, nil
	} else {
		root, _ := users.minter.Mint()
		users.users[uid] = NewUserProfile(uid, root)
		return root, nil
	}
}

func (users UserBase) StoreLSAT(uid UserId, token auth.LSAT) error {
	profile, ok := users.users[uid]

	if ok {
		profile.tokens = append(profile.tokens, token)
		return nil
	} else {
		return errors.New("User not found") // Faudrait avoir une convention d'erreur
	}
}
