package pwd

import (
	"fmt"
	"testing"
)

func TestCheckHash(t *testing.T) {
	hash := HashPwd("12345")
	fmt.Println(hash)
}
func TestHashPwd(t *testing.T) {
	//$2a$04$nxd/7PJlA/7QqVS6fukKTenwUm8fUbUFq4ttdEspwxaWaIZG3cBmK
	ok := CheckHash("$2a$04$nxd/7PJlA/7QqVS6fukKTenwUm8fUbUFq4ttdEspwxaWaIZG3cBmK", "123456")
	fmt.Println(ok)
}
