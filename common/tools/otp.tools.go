package tools

import (
	"crypto/rand"
	"io"

	"grpc-starter/common/constant"
)

// GenerateOTP creates OTP for verification email
func GenerateOTP() (string, string, error) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	otp := make([]byte, constant.Six)
	n, err := io.ReadAtLeast(rand.Reader, otp, constant.Six)
	if n != constant.Six {
		return "", "", err
	}

	for i := 0; i < len(otp); i++ {
		otp[i] = table[int(otp[i])%len(table)]
	}

	otpString := string(otp)

	encrypted, _ := BcryptEncrypt(otpString)

	return otpString, encrypted, nil
}
