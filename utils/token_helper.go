/*
 * @Author: 安俊霖
 * @Date: 2024-11-06 20:50:43
 * @Description:
 */
package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenHelper struct{}

// 签名密钥
var (
	JWTSecret          = "open_kf_2024"
	JWTAlgorithm       = jwt.SigningMethodHS256
	JWTExpirationDelta = 7 * 24 * time.Hour
)

// 生成 JWT
func GenerateToken(userID string) (string, error) {
	payload := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().UTC().Add(JWTExpirationDelta).Unix(),
	}
	token := jwt.NewWithClaims(JWTAlgorithm, payload)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 验证 JWT
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.NewValidationError("token is invalid", jwt.ValidationErrorSignatureInvalid)
}
