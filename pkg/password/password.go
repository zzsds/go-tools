package password

import "golang.org/x/crypto/bcrypt"

// EncryptionPassword 加密
func Encryption(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// VerifyPassword 密码验证
// encryptionPassword 加密串
// verifyPassword 明文密码
// err != nil 检验失败
func Verify(encryptionPassword, verifyPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptionPassword), []byte(verifyPassword)); err != nil {
		return err
	}
	return nil
}
