package auth

import(
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/google/uuid"
	


)
func HashPassword(password string) (string, error){
	byteArray := []byte(password)
	hashPass,err := bcrypt.GenerateFromPassword(byteArray,10)
	if err != nil{
		return "",fmt.Errorf("unable to set password: %w",err)
	}
	return string(hashPass),nil
}

func CheckPasswordHash(hash, password string) error{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil{
		return fmt.Errorf("invalid password %w",err)
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		Subject: userID.String(),
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil{
		return "",err
	}
	return ss,nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(tokenSecret), nil
    })
    if err != nil {
        return uuid.UUID{}, fmt.Errorf("invalid token: %w", err)
    }
    if !token.Valid {
        return uuid.UUID{}, fmt.Errorf("invalid token")
    }
    
    // How can we get the user ID from claims.Subject
	userID,err := uuid.Parse(claims.Subject)
	if err != nil{
		return uuid.UUID{},fmt.Errorf("unable to parse userID: %w",err)
	}
	
    // and convert it to UUID?
	return userID,nil
}