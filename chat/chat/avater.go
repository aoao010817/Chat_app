package main
import (
	"errors"
)

var ErrrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")
// Avatarはユーザーのプロフィールを表す
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}