package test

import (
	"errors"
	"lsat/auth"
	. "lsat/secrets"
)

type UserProfile struct {
	Id     UserId
	Root   Secret
	tokens []auth.Token
}

func NewUserProfile(Id UserId, Root Secret) UserProfile {
	return UserProfile{
		Id,
		Root,
		make([]auth.Token, 10),
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

func (users UserBase) StoreToken(uid UserId, token auth.Token) error {
	profile, ok := users.users[uid]

	if ok {
		profile.tokens = append(profile.tokens, token)
		return nil
	} else {
		return errors.New("User not found") // Faudrait avoir une convention d'erreur
	}
}
