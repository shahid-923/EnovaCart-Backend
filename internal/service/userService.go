package service

import (
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"
	"ecommerce-app/internal/helper"
	"ecommerce-app/internal/repository"

	"ecommerce-app/config"
	"ecommerce-app/pkg/notification"
	"errors"
	"fmt"
	"time"
)

type UserService struct {
	Repo   repository.UserRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s *UserService) Register(input dto.UserSignup) (domain.User, string, error) {
	hashedPassword, err := s.Auth.CreateHashedPassword(input.Password)
	if err != nil {
		return domain.User{}, "", err
	}

	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
	})
	if err != nil {
		return domain.User{}, "", err
	}

	token, err := s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		return domain.User{}, "", err
	}

	return user, token, nil
}

func (s *UserService) Login(email string, password string) (string, error) {
	user, err := s.Repo.FindUser(email)
	if err != nil {
		return "", errors.New("user does not exist")
	}

	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s *UserService) GetProfile(id uint) (domain.User, error) {
	return s.Repo.FindUserByID(id)
}

func (s *UserService) GetUserByID(id uint) (domain.User, error) {
	return s.Repo.FindUserByID(id)
}

func (s *UserService) isVerifiedUser(id uint) bool {

	currentUser, err := s.Repo.FindUserByID(id)
	return err == nil && currentUser.Verified // if true else return false

}
func (s *UserService) GetVerificationCode(e domain.User) error {

	// if user already verified
	if s.isVerifiedUser(e.ID) {
		return errors.New("user already verified")
	}

	// generate verification code
	code, err := s.Auth.GenerateCode()
	if err != nil {
		return err
	}

	// update user
	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}

	_, err = s.Repo.UpdateUser(e.ID, user)

	if err != nil {
		return errors.New("failed to update user")
	}
	user, _ = s.Repo.FindUserByID(e.ID) // get the updated user with code and expiry time

	// send the email to user email
	notificationClient := notification.NewNotificationClient(s.Config)

	message := fmt.Sprintf(
	"Hello,\n\n"+
		"Thank you for registering with MYCOM.\n\n"+
		"Your email verification code is: %d\n\n"+
		"This code is valid for the next 30 minutes.\n\n"+
		"If you did not request this verification, you can safely ignore this email.\n\n"+
		"Regards,\n"+
		"MYCOM Team",
	code,
    )

	err = notificationClient.SendEmail(
		user.Email,
		"Email Verification",
		message,
	)

	if err != nil {
		return err
	}

	// return the verification code
	return nil
}

func (s *UserService) VerifyCode(id uint, code int) error {
	if s.isVerifiedUser(id) {
		return errors.New("user already verified")
	}

	user, err := s.Repo.FindUserByID(id)
	if err != nil {
		return err
	}

	if user.Code != code {
		return errors.New("invalid verification code")
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("verification code expired")
	}

	updatedUser := domain.User{
		Verified: true,
	}

	_, err = s.Repo.UpdateUser(id, updatedUser)
	if err != nil {
		return errors.New("unable to update user")
	}

	return nil
}

func (s *UserService) CreateProfile(id uint, input any) error {
	return nil
}

func (s *UserService) UpdateProfile(id uint, input any) error {
	return nil
}

func (s *UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {

	// find existing user
	user, _ := s.Repo.FindUserByID(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("you have already joined the seller program")
	}

	//update existing user
	seller, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		UserType:  domain.SELLER,
	})

	if err != nil {
		return "", err
	}

	//generate token with new user type
	token, err := s.Auth.GenerateToken(user.ID, user.Email, seller.UserType)

	// create bank account info for the seller

	err = s.Repo.CreateBankAccount(domain.BankAccount{
		BankAccountNumber: input.BankAccountNumber,
		SwiftCode:         input.SwiftCode,
		PaymentType:       input.PaymentType,
		UserID:            id,
	})

	return token, err
}

func (s *UserService) FindCart(id uint) ([]interface{}, error) {
	return nil, nil
}

func (s *UserService) CreateCart(input any, u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s *UserService) CreateOrder(u domain.User) (int, error) {
	return 0, nil
}

func (s *UserService) GetOrders(u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s *UserService) GetOrderById(id uint, uId uint) ([]interface{}, error) {
	return nil, nil
}
