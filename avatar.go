package main

import "errors"

// ErrNoAvatarURLはAvatarインスタンスがアバターのURLを返すことが出来ない
// 場合に発生するエラーです。
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// Avatarはユーザーののプロフィール画像を表す型です、
type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターURLを返します。
	// 問題が発生した場合はエラーを返します。特に、URLを取得出なかった
	// 場合にはErrNoAvatarURLを返します。
	GetAvatarURL(c *client) (string, error)
}
