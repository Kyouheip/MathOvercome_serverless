package middleware

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

// LocalAuthMiddleware はローカル開発時のみ有効なミドルウェア。
// Authorization ヘッダーの JWT からCognito subとusernameを抽出し
// X-User-Sub / X-User-Name ヘッダーにセットする。
// 本番ではAPI Gatewayがこの役割を担うため不要。
func LocalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		parts := strings.Split(strings.TrimPrefix(authHeader, "Bearer "), ".")
		if len(parts) != 3 {
			c.Next()
			return
		}

		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			c.Next()
			return
		}

		var claims map[string]any
		if err := json.Unmarshal(payload, &claims); err != nil {
			c.Next()
			return
		}

		if sub, ok := claims["sub"].(string); ok {
			c.Request.Header.Set("X-User-Sub", sub)
		}
		if name, ok := claims["name"].(string); ok {
			c.Request.Header.Set("X-User-Name", name)
		}

		c.Next()
	}
}
