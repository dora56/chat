package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testCharUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testCharUser)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合, AuthAvatar.GetAvatarURLは ErrorNotAvatarURLを使うべきです。")
	}
	// 値をセットします
	testUrl := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testCharUser.User = testUser
	testUser.On("AvatarURL").Retrurn(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testCharUser)
	if err != nil {
		t.Error("値が存在する場合、authAvatar.GetAvatarURLはエラーを返すべきではありません")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GavatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvitar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	//テスト用のアバターのファイルを生成します。
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte(), 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("fileSystemAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "/avatar/abc.jpg" {
		t.Error("fileSystemAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}
