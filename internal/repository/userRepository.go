package repository

import (
	"ecommerce-app/internal/domain"
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(user domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserByID(id uint) (domain.User, error)
	UpdateUser(id uint, user domain.User) (domain.User, error)
	CreateBankAccount(e domain.BankAccount) error

	// Cart
	FindCartItems(uId uint) ([]domain.Cart, error)
	FindCartItem(uId uint, pId uint) (domain.Cart, error)
	CreateCart(c domain.Cart) error
	UpdateCart(c domain.Cart) error
	DeleteCartById(id uint) error
	DeleteCartItems(uId uint) error

	// Order
	CreateOrder(o domain.Order) error
	FindOrders(uId uint) ([]domain.Order, error)
	FindOrderById(id uint, uId uint) (domain.Order, error)

	// Profile
	CreateProfile(e domain.Address) error
	UpdateProfile(e domain.Address) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateBankAccount(e domain.BankAccount) error {
	return r.db.Create(&e).Error
}

func (r *userRepository) CreateUser(user domain.User) (domain.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		log.Printf("error creating user: %v\n", err)
		return domain.User{}, errors.New("failed to create user")
	}
	return user, nil
}

func (r *userRepository) FindUser(email string) (domain.User, error) {

	var user domain.User

	err := r.db.Preload("Address").First(&user, "email = ?", email).Error
	if err != nil {
		log.Printf("error finding user: %v\n", err)
		return domain.User{}, errors.New("user doesn't exist")
	}

	return user, nil
}

func (r *userRepository) FindUserByID(id uint) (domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").Preload("Cart").Preload("Orders").First(&user, "id = ?", id).Error // preload val same as fields in user.go
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {
	var user domain.User

	err := r.db.Model(&domain.User{}).
		Where("id = ?", id).
		Updates(u).
		First(&user).Error

	if err != nil {
		log.Printf("error updating user: %v\n", err)
		return domain.User{}, errors.New("failed to update user")
	}

	return user, nil
}

func (r *userRepository) FindCartItems(uId uint) ([]domain.Cart, error) {
	var carts []domain.Cart

	err := r.db.
		Where("user_id = ?", uId).Find(&carts).Error

	if err != nil {
		return nil, err
	}

	return carts, nil
}

func (r *userRepository) FindCartItem(uId uint, pId uint) (domain.Cart, error) {
	var cartItem domain.Cart
	err := r.db.Where("user_id = ? AND product_id = ?", uId, pId).First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Cart{}, nil
		}
		return domain.Cart{}, err
	}
	return cartItem, nil
}

func (r *userRepository) CreateCart(c domain.Cart) error {
	return r.db.Create(&c).Error
}

func (r *userRepository) UpdateCart(c domain.Cart) error {

	var cart domain.Cart
	err := r.db.Model(&cart).Clauses(clause.Returning{}).Where("id = ?", c.ID).Updates(c).Error
	return err
}

func (r *userRepository) DeleteCartById(id uint) error {

	err := r.db.Delete(&domain.Cart{}, id).Error
	return err
}

func (r *userRepository) DeleteCartItems(uId uint) error {
	err := r.db.Where("user_id = ?", uId).Delete(&domain.Cart{}).Error
	return err
}

func (r *userRepository) CreateProfile(e domain.Address) error {
	err := r.db.Create(&e).Error

	if err != nil {
		log.Printf("error on creating profile with address %v", err)
		return errors.New("failed to create profile")
	}
	return nil
}

func (r *userRepository) UpdateProfile(e domain.Address) error {

	err := r.db.Where("user_id=?", e.UserId).Updates(e).Error
	if err != nil {
		log.Printf("error on updating profile with address %v", err)
		return errors.New("failed to update profile")
	}
	return nil
}

func (r *userRepository) CreateOrder(o domain.Order) error {

	err := r.db.Create(&o).Error

	if err != nil {
		log.Printf("error on creating order %v", err)
		return errors.New("failed to create order in database")
	}
	return nil
}

func (r *userRepository) FindOrders(uId uint) ([]domain.Order, error) {

	var orders []domain.Order
	err := r.db.Where("user_id", uId).Find(&orders).Error

	if err != nil {
		log.Printf("error on fetching orders %v", err)
		return nil, errors.New("failed to fetch orders")
	}
	return orders, nil
}

func (r *userRepository) FindOrderById(id uint, uId uint) (domain.Order, error) {

	var order domain.Order
	err := r.db.Preload("Items").First(&order, id).Error

	if err != nil {
		log.Printf("error on fetching orders %v", err)
		return domain.Order{}, errors.New("failed to fetch orders")
	}
	return order, nil
}
