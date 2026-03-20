package user

import (
	"os"
	"os/user"
)

func GetOsUsername() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	if s := os.Getenv("USER"); s != "" {
		return s
	}
	if s := os.Getenv("USERNAME"); s != "" {
		return s
	}
	return "user"
}
