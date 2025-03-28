/* Copyright 2024 GuanceCloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
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
