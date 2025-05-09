package usecases

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"server2/application/entities"
	"server2/application/infrastructure"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

const ()

type PassowrdHashUseCase struct {
	times           uint32
	availableMemomy uint32
	threads         uint8
	saltLen         uint32
}

func NewPassowrdHashUseCase() PassowrdHashUseCase {
	return PassowrdHashUseCase{
		times:           3,
		availableMemomy: 64 * 1024,
		threads:         2,
		saltLen:         32,
	}
}

// rawstring -> argon2 -> params + salt(base64) + hash(base64)
func (hp *PassowrdHashUseCase) Hash(password string) []byte {
	salt := make([]byte, hp.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	fmt.Println(hp.times)
	hash := argon2.IDKey([]byte(password), salt, hp.times, hp.availableMemomy, hp.threads, hp.saltLen)
	// fmt.Println(salt)
	hashB64, saltB64 := []byte(base64.RawStdEncoding.EncodeToString(hash)), []byte(base64.RawStdEncoding.EncodeToString(salt))

	hashString := fmt.Sprintf("%d$%d$%d$%d$%s$%s",
		hp.times, hp.availableMemomy, hp.threads, hp.saltLen,
		saltB64,
		hashB64,
	)

	return []byte(hashString)
}

func (hp *PassowrdHashUseCase) parseParams(timeStr, memoryStr, threadsStr, keyLenStr string) (uint32, uint32, uint8, uint32) {
	var time, memory, keyLen uint32
	var threads uint8
	fmt.Sscanf(timeStr, "%d", &time)
	fmt.Sscanf(memoryStr, "%d", &memory)
	fmt.Sscanf(threadsStr, "%d", &threads)
	fmt.Sscanf(keyLenStr, "%d", &keyLen)
	return time, memory, threads, keyLen
}

// process and compare the password hash with the input hash
func (hp *PassowrdHashUseCase) VerifyPassword(password string, hash string) bool {

	// time$availableMemomy$threads$saltLen$salt(b64)$hash(b64)
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}
	time, memory, threads, keyLen := hp.parseParams(parts[0], parts[1], parts[2], parts[3])
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		panic(err)
	}
	hashStored, err := base64.RawStdEncoding.DecodeString(parts[5])

	if err != nil {
		fmt.Println(err)
		return false
	}

	hashComputed := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	return subtle.ConstantTimeCompare(hashStored, hashComputed) == 1
}

type JwtUseCase struct {
	sessionSecret string
}

func NewJWTUseCase(sessionSecret string) *JwtUseCase {
	return &JwtUseCase{sessionSecret: sessionSecret}
}

// generate a JWT token with userid and string as subject ('sub' claim)
func (j *JwtUseCase) GenerateJWT(userId string, address string) (string, error) {
	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%s$%s", userId, address),
		"exp": time.Now().Add(time.Second * 3600).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.sessionSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// returns userid and address without port (string,string)
func (j *JwtUseCase) ParseJWTSubject(subject string) (string, string) {
	splited := strings.Split(subject, "$")
	return splited[0], strings.Split(splited[1], ":")[0]
}

// verify if token signing method is HMAC
func (j *JwtUseCase) verifySigningMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
	}
	return []byte(j.sessionSecret), nil
}

// verify if token is valid and return the object
func (j *JwtUseCase) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, j.verifySigningMethod)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// verify if authorization header (bearer xxxxxxxxxx) is a valid jwt, and if jwt address is compatible with requester ip
func (j *JwtUseCase) TokenFromBearer(bearer string) (*jwt.Token, error) {
	authorization, _ := strings.CutPrefix(bearer, "Bearer ")
	token, err := j.ValidateJWT(string(authorization))

	if err != nil {
		return nil, err
	}

	// subject, err := token.Claims.GetSubject()

	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }

	// _, userAddress := j.ParseJWTSubject(subject)

	// if userAddress != strings.Split(r.RemoteAddr, ":")[0] {
	// 	return nil, fmt.Errorf("jwt address does not match with request address")
	// }

	return token, nil
}

type GetUserFromJWTBearerUseCase struct {
	jwtUseCase *JwtUseCase
	userRepo   infrastructure.UserRepository
}

func NewGetUserFromJWTBearerUseCase(userRepo infrastructure.UserRepository, jwtUseCase *JwtUseCase) *GetUserFromJWTBearerUseCase {
	return &GetUserFromJWTBearerUseCase{userRepo: userRepo, jwtUseCase: jwtUseCase}
}

// get the user with userID stored at jwt subject
func (gu *GetUserFromJWTBearerUseCase) Execute(bearer string) (*entities.User, error) {
	token, err := gu.jwtUseCase.TokenFromBearer(bearer)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	jwtSubject, ok := token.Claims.(jwt.MapClaims)["sub"].(string)

	if !ok || jwtSubject == "" {
		return nil, fmt.Errorf("invalid subject")
	}

	userid, _ := gu.jwtUseCase.ParseJWTSubject(jwtSubject)
	userUseCase := NewUserUseCase(gu.userRepo)
	requesterUser, err := userUseCase.Get(userid)
	if err != nil {
		return nil, err
	}
	return requesterUser, nil
}

type UserUseCase struct {
	repo infrastructure.UserRepository
}

func NewUserUseCase(repo infrastructure.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) Save(username string, password string, roleID string, secret string) (string, error) {
	fmt.Println("saving user", username)
	id, err := u.repo.Save(username, password, roleID, secret)
	return id, err
}

func (u *UserUseCase) Update(id, username string, password string, roleID string, secret string) error {
	fmt.Println("updating user", id)
	err := u.repo.Update(id, username, password, roleID, secret)
	return err
}

func (u *UserUseCase) Get(id string) (*entities.User, error) {
	fmt.Println("getting user", id)
	user, err := u.repo.Get(id)
	return user, err
}

func (u *UserUseCase) Delete(id string) error {
	fmt.Println("deleting user", id)
	err := u.repo.Delete(id)
	return err
}

func (u *UserUseCase) SearchByName(regex string) ([]*entities.User, error) {
	fmt.Println("searching user by name", regex)
	users, err := u.repo.SearchByName(regex)
	return users, err
}

func (u *UserUseCase) SearchByRoleID(roleName string) ([]*entities.User, error) {
	fmt.Println("searching user by role", roleName)
	user, err := u.repo.SearchByRoleID(roleName)
	return user, err
}
