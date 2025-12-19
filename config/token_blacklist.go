package config

import "sync"

// TokenBlacklist menyimpan token JWT yang sudah logout
var TokenBlacklist sync.Map

func BlacklistToken(token string) {
	TokenBlacklist.Store(token, true)
}

func IsTokenBlacklisted(token string) bool {
	_, found := TokenBlacklist.Load(token)
	return found
}
