package auth

import(
	"golang.org/x/crypto/bcrypt"
	"fmt"
)


func HashPassword(password string) (string, error){
	byteArray := []byte(password)
	hashPass,err := bcrypt.GenerateFromPassword(byteArray,10)
	if err != nil{
		return fmt.Errorf("unable to set password: %w",err)
	}
	return string(hashPass),nil
}

func CheckPasswordHash(password, hash string) error{
	err := bcrypt.CompareHashAndPassword(password, []byte(hash))
	if err != nil{
		return fmt.Errorf("invalid password %w",err)
	}
	return nil
}