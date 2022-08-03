package jwt

import (
	"testing"
	"time"

	"github.com/zzsds/go-tools/auth"
)

var (
	// RSA2 密钥 base64 加密所得参数
	private = `LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWZ0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVVBd2dnRThBZ0VBQWtFQXZ6WnJjakxKOXFKcURhZ00KeGh5bGgrOFcxUkVtb2NmdjF0WWFmS1Z3WlNsTy9WWkduOER4WXgvNlBWMUZJdjBTaE9pN0NWSWpMNE9sM0lrZQo0ZjFCNHdJREFRQUJBa0VBcGZPbU54dkxXeW5FbjR1ZFlvZlVSbkVFUVBHOHRLWmhDdlVSVWVNSDlGTHdWelhkCjc0TnpHL0s5bzduNXlHOWhFdDl0ZnpabkZGeXVHanNxWUFFbDJRSWhBT2M5a3JHS0ZnZTcwYi9PdDdKOGpTMGwKaGpOUVJQQ0paeDhSRUhjU2FFZk5BaUVBMDYrb3hVSmhsQ01Lc1AxZVNZUGNnYThDcmdvbHFJVmJGdjdySmU5KwpvRzhDSVFERkQ1bjFwc0hEY1hIOFRZUUtuVTRLVFZJaVpLTjdnUHphWXNadlVzWi9lUUlnZlBOR3o0anJTQ0dYCkFmbk1XZUIzbkNUTmxDVnhMUlBxUEp5aitIUnhiZ2tDSVFDbkU3cGNJazNJbi9GcDNUVHd5VDZRMFJ2aGhPWmsKWFFsR2FNRWFZQmNjeXc9PQotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0t`

	// RSA2 密钥 base64 加密所得参数
	pub = `LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTDgyYTNJeXlmYWlhZzJvRE1ZY3BZZnZGdFVSSnFISAo3OWJXR255bGNHVXBUdjFXUnAvQThXTWYrajFkUlNMOUVvVG91d2xTSXkrRHBkeUpIdUg5UWVNQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==`

	tok = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiRVhDSEFOR0UiLCJzY29wZXMiOm51bGwsIm1ldGFkYXRhIjp7Im1vYmlsZSI6IjE1NzM2NDUwOTI1In0sImV4cCI6MTYxNzE3MzQ3OSwiaXNzIjoiZ28ua3JhdG9zIiwic3ViIjoiNCJ9.snlxdhlmle-KWgfJENRoPsSBMYKJw3Qu-J\nYlho8hM6O5kiPyz-4Xyx9i7vGK9GCogG6U3L1hPxmj5crvWuZeeI2SsqHYXofXpNFh8iFiaMjfoAj0oSAHEKDm3VkqeJo-zyHJEQ9jRATtMgSfxTU-QDIAIyy8iOTyb_tIXoFU_l-Fliq0vlk8T4yI32MeBx-4JvO6ffaBNNywx6ZA0YoP70dDcpFiz1-MpbkHmwtDPDli0vBiTaLfU2h0ceT746gs4oX1AHmfA-\nOOz1O6yPeLquRaRT-GsZQjNb3zY7lWGX-96bkR2FUL8ZDs2FsYeamxHFNoiLobBP1YL8vdK_FSEk3rtm8lVeK3ZBpQ5eS5TCjHT6lA7qZeXpSJvbPhMpsi9DcjJzxxJQG7G8vKbbox7r8IouXv3-IYrAHMBujJuXIsZv5EgZWvuojNcC-dIZAnY8hsFq3lpUJaGMQLtoqGCa647Yp-w1VlJpR0Nkxq2j9hPbCBgT\n0PIG0USSCQ1fdFQBy5x-DrhTGaNqtcHN0mFwo-H-6OdYoCOhFAFr3iKHMP9qxctgzh7EntmpDdIm7dw5TlGYiVGrD6H47jxJxOGu50tPhTnwz9OmrJIIP2AkZcr1gSu7MdwItlpFf5P_sPI6Am6xV-uhDeVXP45fyx7vA-QdTIv5PJyjh7_0m_V4I"
)

func TestNewAuth(t *testing.T) {
	a := NewAuth(func(o *auth.Options) {
		o.PrivateKey = private
		o.PublicKey = pub
	})
	account, err := a.Generate("1", func(o *auth.GenerateOptions) {
		o.Scopes = []string{"1"}
		o.Type = "EXCHANGE"
		o.Metadata = map[string]string{"mobile": "12564586521"}
	})
	if err != nil {
		t.Fatal(err, 12321)
	}
	t.Log(account)
	to, err := a.Token(func(o *auth.TokenOptions) {
		o.Expiry = 24 * time.Hour
		o.Secret = account.Secret
	})
	if err != nil {
		t.Fatal(err)
	}
	account, err = a.Inspect(to.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(account)
}
