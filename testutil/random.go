package testutil

import (
	"math/rand"
	"restaurant-order-management/models"
	"strings"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"
var alphaNumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandomString(n int) string {
	k := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomRole() (role models.Role) {
	role = models.Role{
		Name: RandomString(6),
	}
	return
}

func RandomUser() (user models.User) {
	user = models.User{
		Name:     RandomString(6),
		Email:    RandomString(6),
		Password: RandomString(3),
	}
	return
}
