package service

import (
	"ecommerce-app/internal/domain"
	"ecommerce-app/internal/dto"
	"ecommerce-app/internal/helper"
	"ecommerce-app/internal/repository"
	"log"

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
	CRepo  repository.CatalogueRepository
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
			"Thank you for registering with NovaCart.\n\n"+
			"Your email verification code is: %d\n\n"+
			"This code is valid for the next 30 minutes.\n\n"+
			"If you did not request this verification, you can safely ignore this email.\n\n"+
			"Regards,\n"+
			"NovaCart Team",
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

func (s *UserService) CreateProfile(id uint, input dto.ProfileInput) error {

	// update user
	_, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})

	if err != nil {
		return err
	}

	// create address
	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserId:       id,
	}

	err = s.Repo.CreateProfile(address)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetProfile(id uint) (*domain.User, error) {

	user, err := s.Repo.FindUserByID(id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateProfile(id uint, input dto.ProfileInput) error {

	user, err := s.Repo.FindUserByID(id)

	if err != nil {
		return err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}

	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = s.Repo.UpdateUser(id, user)

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserId:       id,
	}

	err = s.Repo.UpdateProfile(address)
	if err != nil {
		return err
	}

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

func (s *UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {

	if input.ProductId == 0 {
		return nil, errors.New("please provide the valid product id")
	}

	// check if cart item already exists for this user+product
	cart, err := s.Repo.FindCartItem(u.ID, input.ProductId)
	if err != nil {
		return nil, errors.New("error checking cart item")
	}

	switch {
	case cart.ID > 0 && input.Qty < 1:
		// existing item, qty < 1 -> remove it
		if err := s.Repo.DeleteCartById(cart.ID); err != nil {
			log.Printf("Error deleting cart item %v", err)
			return nil, errors.New("error on deleting cart item")
		}

	case cart.ID > 0 && input.Qty >= 1:
		// existing item, qty given -> update it
		cart.Qty = input.Qty
		if err := s.Repo.UpdateCart(cart); err != nil {
			return nil, errors.New("error on updating cart item")
		}

	case cart.ID == 0 && input.Qty >= 1:
		// no existing item -> create new one
		product, err := s.CRepo.FindProductById(int(input.ProductId))
		if err != nil {
			return nil, errors.New("product not found to create cart item")
		}

		newCart := domain.Cart{
			ProductId: input.ProductId,
			UserId:    u.ID,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Qty:       input.Qty,
			Price:     int64(product.Price),
			SellerId:  uint(product.UserId),
		}

		if err := s.Repo.CreateCart(newCart); err != nil {
			return nil, errors.New("error on creating cart item")
		}
	}

	return s.Repo.FindCartItems(u.ID) // return the new total cart to save api calls
}

func (s *UserService) FindCart(id uint) ([]domain.Cart, error) {

	return s.Repo.FindCartItems(id)
}

func (s *UserService) CreateOrder(u domain.User) (int, error) {

	// find the cart items for the user
	cartItems, err := s.Repo.FindCartItems(u.ID)

	if err != nil {
		return 0, errors.New("error on finding cart items")
	}

	if len(cartItems) == 0 {
		return 0, errors.New("cart is empty cannot create the order")
	}

	// find success payment reference status
	paymentId := ""     // get from payment providers
	transactionId := "" // expose to the user

	orderRef, _ := helper.RandomNumbers(8)

	// create order with generated OrderNo
	var amount float64
	var orderItems []domain.OrderItem

	for _, item := range cartItems {

		amount += float64(item.Price) * float64(item.Qty)
		orderItems = append(orderItems, domain.OrderItem{

			ProductId: item.ProductId,
			Qty:       item.Qty,
			Price:     float64(item.Price),
			Name:      item.Name,
			ImageUrl:  item.ImageUrl,
			SellerId:  item.SellerId,
		})

	}

	order := domain.Order{
		UserId:         u.ID,
		PaymentId:      paymentId,
		TransactionId:  transactionId,
		OrderRefNumber: uint(orderRef),
		Amount:         amount,
		Items:          orderItems,
	}

	err = s.Repo.CreateOrder(order)

	if err != nil {
		return 0, err
	}

	// send email to the user with order details
	// send email to the seller with order details

	// remove the cart items after processing above
	err = s.Repo.DeleteCartItems(u.ID)
	log.Printf("Error deleting cart items %v")

	// return orderRef no

	return orderRef, nil
}

func (s *UserService) GetOrders(u domain.User) ([]domain.Order, error) {

	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, errors.New("error fetching orders")
	}
	return orders, nil
}

func (s *UserService) GetOrderById(id uint, uId uint) (domain.Order, error) {

	order, err := s.Repo.FindOrderById(id, uId)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
