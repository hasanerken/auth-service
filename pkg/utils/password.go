package encryption

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPwd(pwd string) []byte {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Warning("Cannot hash the password", err)
	}
	return hashedPwd
}

func ComparePwd(hashedPwdFromDB []byte, pwd string) (bool, error) {
	fmt.Printf("pwd %v\n", pwd)
	fmt.Printf("pwdfromdb %v\n", hashedPwdFromDB)
	if err := bcrypt.CompareHashAndPassword(hashedPwdFromDB, []byte(pwd)); err != nil {
		log.Warning("Can not compare", err)
		return false, err
	}
	return true, nil
}
