package main
import  "testing"
func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返します。")
	}
	testUrl := "http://url-to-avatar"
	client.userData = map[string]interface{}{"avatar_url": testUrl} 
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.error("AuthAvatar.GetAvatarURLは正しいURLを返すべきです")
	}
}