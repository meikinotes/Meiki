package auth

import (
	"context"
	"errors"
	"regexp"

	"github.com/RainingComputers/Meiki/log"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	token_coll *mongo.Collection
	user_coll  *mongo.Collection
}

type User struct {
	Username     string `bson:"username"`
	PasswordHash []byte `bson:"password"`
}

type UserTokens struct {
	Username string   `bson:"username"`
	Tokens   []string `bson:"tokens"`
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrMissingUser       = errors.New("unable to find user in DB")
	ErrUnableToLogOut    = errors.New("unable to log out this user")
	ErrPasswordMismatch  = errors.New("could not login user due to password mismatch")
	ErrMissingUserTokens = errors.New("unable to find user tokens in DB")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidPassword   = errors.New("invalid password")
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+([_-]?[a-zA-Z0-9])*$`)

func getToken() string {
	return uuid.NewString()
}

func CreateAuth(ctx context.Context, token_coll *mongo.Collection, user_coll *mongo.Collection) (Auth, error) {
	mod := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := user_coll.Indexes().CreateOne(ctx, mod)

	if err != nil {
		log.Error("unable to create unique index in user collection", zap.Error(err))
		return Auth{}, err
	}

	_, err = token_coll.Indexes().CreateOne(ctx, mod)

	if err != nil {
		log.Error("unable to create unique index in token collection", zap.Error(err))
		return Auth{}, err
	}

	return Auth{token_coll, user_coll}, nil
}

func validateCredentials(username string, password string) error {
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}

	if len(password) < 5 {
		return ErrInvalidPassword
	}

	return nil
}

func (a Auth) storeCredentialsInDB(ctx context.Context, user User) error {

	_, err := a.user_coll.InsertOne(ctx, user)

	if mongo.IsDuplicateKeyError(err) {
		log.Error("user already exists", zap.Error(err))
		return ErrUserAlreadyExists
	}

	if err != nil {
		log.Error("could not store credentials in DB", zap.Error(err))
		return err
	}

	return nil
}

func (a Auth) Create(ctx context.Context, username string, password string) error {
	err := validateCredentials(username, password)

	if err != nil {
		log.Info("abort creating user because invalid credentials", zap.Error(err))
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	if err != nil {
		log.Fatal("could not create password hash for user", zap.Error(err))
		return err
	}

	err = a.storeCredentialsInDB(ctx, User{
		username,
		passwordHash,
	})

	if err != nil {
		log.Error("could not store credentials while creating user", zap.Error(err))
		return err
	}

	return nil
}

func (a Auth) deleteUserInDB(ctx context.Context, username string) error {
	result, err := a.user_coll.DeleteOne(ctx, bson.M{"username": username})

	if err != nil {
		log.Error("unable to delete user from DB", zap.Error(err))
		return err
	}

	if result.DeletedCount == 0 {
		log.Error("could not find user in DB")
		return ErrMissingUser
	}

	return nil
}

func (a Auth) deleteUserTokensInDB(ctx context.Context, username string) error {
	result, err := a.token_coll.DeleteOne(ctx, bson.M{"username": username})

	if err != nil {
		log.Error("unable to delete user token from DB", zap.Error(err))
		return err
	}

	if result.DeletedCount == 0 {
		log.Error("could not find user token in DB")
		return ErrMissingUserTokens
	}

	return nil
}

func (a Auth) Delete(ctx context.Context, username string) error {
	err := a.deleteUserInDB(ctx, username)

	if err != nil {
		log.Error("unable to delete user", zap.Error(err))
		return err
	}

	err = a.deleteUserTokensInDB(ctx, username)

	if err == ErrMissingUserTokens {
		log.Warn("new user deleted without ever logging in")
		return nil
	}

	if err != nil {
		log.Error("unable to delete user tokens", zap.Error(err))
		return err
	}

	return nil
}

func (a Auth) getPasswordHashFromDB(ctx context.Context, username string) ([]byte, error) {
	var user User

	err := a.user_coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Info("password missmatch for user", zap.Error(err))
		return nil, ErrMissingUser
	}

	if err != nil {
		log.Error("unable to read password for user", zap.Error(err))
		return nil, err
	}

	return user.PasswordHash, nil
}

func (a Auth) PasswordMatches(ctx context.Context, username, password string) (bool, error) {
	passwordHash, err := a.getPasswordHashFromDB(ctx, username)

	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))

	if err != nil {
		log.Info("password mismatch for user", zap.String("username", username))
		return false, nil
	}

	return true, nil
}

func (a Auth) CreateToken(ctx context.Context, username string) (string, error) {
	newToken := getToken()

	result, err := a.token_coll.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$push": bson.M{"tokens": newToken}})

	if err != nil {
		log.Error("unable to update token into DB", zap.Error(err))
		return "", err
	}

	if result.MatchedCount == 0 {
		newUserTokens := UserTokens{
			Username: username,
			Tokens:   []string{newToken},
		}

		_, err := a.token_coll.InsertOne(ctx, newUserTokens)

		if err != nil {
			log.Error("unable to insert token into DB", zap.Error(err))
			return "", err
		}
	}

	return newToken, nil
}

func (a Auth) Login(ctx context.Context, username string, password string) (string, error) {
	err := validateCredentials(username, password)

	if err != nil {
		log.Info("abort login because invalid credentials", zap.Error(err))
		return "", err
	}

	match, err := a.PasswordMatches(ctx, username, password)

	if err != nil {
		return "", err
	}

	if !match {
		return "", ErrPasswordMismatch
	}

	token, err := a.CreateToken(ctx, username)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a Auth) ReadTokensFromDB(ctx context.Context, username string) ([]string, error) {
	var userTokens UserTokens
	err := a.token_coll.FindOne(ctx, bson.M{"username": username}).Decode(&userTokens)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrMissingUserTokens
	}

	if err != nil {
		return nil, err
	}

	return userTokens.Tokens, nil
}

func (a Auth) Authenticate(ctx context.Context, username string, token string) (bool, error) {

	existingTokens, err := a.ReadTokensFromDB(ctx, username)

	if err != nil {
		return false, err
	}

	for _, t := range existingTokens {
		if t == token {
			return true, nil
		}
	}

	return false, nil
}

func (a Auth) deleteSingleTokenFromDB(ctx context.Context, username string, token string) error {
	result, err := a.token_coll.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$pull": bson.M{"tokens": token}})

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		log.Error("unable to find existing token", zap.Error(err))
		return ErrMissingUserTokens
	}

	return nil
}

func (a Auth) Logout(ctx context.Context, username string, token string) error {
	err := a.deleteSingleTokenFromDB(ctx, username, token)

	if err != nil {
		return err
	}

	log.Info("logged out user successfully", zap.String("username", username))

	return nil
}
