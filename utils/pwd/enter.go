package pwd

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPwd(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}
func CheckHash(hashpwd string, pwd string) bool {
	bytehash := []byte(hashpwd)
	err := bcrypt.CompareHashAndPassword(bytehash, []byte(pwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
