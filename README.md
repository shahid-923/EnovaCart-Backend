# 🛒 EnovaCart — E-Commerce Application

EnovaCart is a full-stack e-commerce application with a **Golang backend** providing REST APIs for user authentication, email verification, profiles, product catalogue management, shopping carts, orders, seller management, and Stripe payment processing.

The backend is designed using a layered architecture with separate **handlers, services, repositories, domain models, DTOs, helpers, and external service integrations**.

---

# 📌 Features

* User registration and login
* JWT-based authentication
* Password hashing with bcrypt
* Email verification using 6-digit OTP
* OTP expiration
* User profile management
* Product catalogue
* Category management
* Seller registration
* Seller-specific authorization
* Product management for sellers
* Shopping cart management
* Order creation and order history
* Seller order management
* Stripe Checkout integration
* Stripe payment-session status checking
* Stripe payment-session expiration
* Stripe metadata for `order_id` and `user_id`
* PostgreSQL database
* GORM ORM
* Docker support
* GitHub Actions CI
* Air hot reload for development

---

# 🏗️ Architecture

```text
                         ┌───────────────────┐
                         │      Client       │
                         │ React / Postman   │
                         └─────────┬─────────┘
                                   │
                                   ▼
                         ┌───────────────────┐
                         │    REST API      │
                         │      Fiber       │
                         └─────────┬─────────┘
                                   │
                                   ▼
                         ┌───────────────────┐
                         │     Handlers      │
                         └─────────┬─────────┘
                                   │
                                   ▼
                         ┌───────────────────┐
                         │     Services      │
                         │  Business Logic  │
                         └─────────┬─────────┘
                                   │
                    ┌──────────────┴──────────────┐
                    │                             │
                    ▼                             ▼
           ┌─────────────────┐          ┌─────────────────┐
           │   Repository    │          │     Helpers     │
           │   Data Access   │          │ JWT / OTP etc.  │
           └────────┬────────┘          └─────────────────┘
                    │
                    ▼
           ┌─────────────────┐
           │   PostgreSQL    │
           └─────────────────┘

                    External Services
                    ─────────────────

                  ┌──────────────┐
                  │    Stripe    │
                  │   Payments   │
                  └──────────────┘

                  ┌──────────────┐
                  │ Notification │
                  │    / Email    │
                  └──────────────┘
```

---

# 🧰 Tech Stack

| Category         | Technology           |
| ---------------- | -------------------- |
| Backend          | Go                   |
| Web Framework    | Fiber                |
| ORM              | GORM                 |
| Database         | PostgreSQL           |
| Authentication   | JWT                  |
| Password Hashing | bcrypt               |
| Payment Gateway  | Stripe               |
| Email            | Notification Service |
| API Testing      | Postman              |
| Containerization | Docker               |
| CI               | GitHub Actions       |
| Hot Reload       | Air                  |
| Version Control  | Git / GitHub         |
| Cloud Knowledge  | AWS / GCP            |

---

# 📁 Project Structure

```text
ecommerce-app/
│
├── config/
│
├── infra/
│
├── internal/
│   │
│   ├── api/
│   │   ├── graphql/
│   │   ├── grpc/
│   │   └── rest/
│   │       └── handlers/
│   │
│   ├── domain/
│   │
│   ├── dto/
│   │
│   ├── helper/
│   │
│   ├── repository/
│   │
│   └── service/
│
├── pkg/
│   ├── notification/
│   └── payment/
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── docker-compose.yml
├── go.mod
└── README.md
```

---

# 🔐 Authentication

Protected endpoints require the JWT returned during login or registration.

Use the following header:

```http
Authorization: Bearer <JWT_TOKEN>
```

The JWT contains information such as:

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "user",
  "exp": 1780000000
}
```

Seller-only endpoints require a JWT containing the seller role.

---

# 👤 User APIs

Base URL:

```text
http://localhost:9000/users
```

---

## 1. Register User

### Endpoint

```http
POST /users/register
```

### Authentication

Public

### Request

```json
{
  "email": "user@example.com",
  "password": "test1234",
  "phone": "9876543210"
}
```

### Success Response

```json
{
  "message": "user created successfully",
  "token": "<JWT_TOKEN>"
}
```

### Status

```text
201 Created
```

---

## 2. Login

### Endpoint

```http
POST /users/login
```

### Authentication

Public

### Request

```json
{
  "email": "user@example.com",
  "password": "test1234"
}
```

### Success Response

```json
{
  "message": "login successful",
  "token": "<JWT_TOKEN>"
}
```

### Status

```text
200 OK
```

### Invalid Credentials

```json
{
  "message": "invalid credentials"
}
```

### Status

```text
401 Unauthorized
```

---

# 📧 Email Verification APIs

All verification endpoints require authentication.

---

## 3. Generate Verification OTP

### Endpoint

```http
GET /users/verify
```

### Authentication

Required

```http
Authorization: Bearer <JWT_TOKEN>
```

### Request Body

None.

### Success Response

```json
{
  "message": "Verification email sent"
}
```

### OTP

The application generates a **6-digit OTP**.

Example:

```text
483921
```

The OTP is stored with an expiration time of **30 minutes**.

---

## 4. Verify OTP

### Endpoint

```http
POST /users/verify
```

### Authentication

Required

### Request

```json
{
  "code": 483921
}
```

### Success Response

```json
{
  "message": "code verified"
}
```

### Invalid OTP

```json
{
  "message": "invalid verification code"
}
```

### Expired OTP

```json
{
  "message": "verification code expired"
}
```

---

# 👨‍💼 Profile APIs

---

## 5. Create Profile

### Endpoint

```http
POST /users/profile
```

### Authentication

Required

### Request

```json
{
  "first_name": "Md",
  "last_name": "Ansari",
  "address_input": {
    "address_line1": "123 Main Street",
    "address_line2": "Apartment 2",
    "city": "Bengaluru",
    "country": "India",
    "post_code": "560001"
  }
}
```

### Success Response

```json
{
  "message": "Profile created successfully"
}
```

---

## 6. Get Profile

### Endpoint

```http
GET /users/profile
```

### Authentication

Required

### Request Body

None.

### Success Response

The response contains the user's profile information.

Example structure:

```json
{
  "message": "Profile fetched",
  "profile": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "Md",
    "last_name": "Ansari"
  }
}
```

---

## 7. Update Profile

### Endpoint

```http
PATCH /users/profile
```

### Authentication

Required

### Request

```json
{
  "first_name": "Md",
  "last_name": "Shahid",
  "address_input": {
    "address_line1": "Updated Address",
    "address_line2": "",
    "city": "Bengaluru",
    "country": "India",
    "post_code": "560001"
  }
}
```

### Success Response

```json
{
  "message": "Profile updated successfully"
}
```

---

# 🛍️ Product Catalogue APIs

## Public Product APIs

---

## 8. Get Products

### Endpoint

```http
GET /products
```

### Authentication

Public

### Request Body

None.

### Success Response

```json
{
  "message": "products fetched successfully",
  "products": []
}
```

The exact product response fields depend on the `domain.Product` model used by the application.

---

## 9. Get Product By ID

### Endpoint

```http
GET /products/:id
```

### Example

```http
GET /products/1
```

### Authentication

Public

### Request Body

None.

### Success Response

```json
{
  "message": "product fetched successfully",
  "product": {}
}
```

---

# 🏷️ Seller Category APIs

These APIs require seller authorization.

```http
Authorization: Bearer <SELLER_JWT_TOKEN>
```

---

## 10. Get Categories

```http
GET /seller/categories
```

### Success Response

```json
{
  "message": "categories fetched successfully",
  "categories": []
}
```

---

## 11. Get Category By ID

```http
GET /seller/categories/:id
```

Example:

```http
GET /seller/categories/1
```

### Success Response

```json
{
  "message": "category fetched successfully",
  "category": {}
}
```

---

## 12. Create Category

```http
POST /seller/categories
```

### Authentication

Seller required.

### Request

The exact category DTO fields were not available in the current project context, so they should be taken directly from the project's category DTO.

### Success Response

```json
{
  "message": "category created successfully"
}
```

---

## 13. Update Category

```http
PATCH /seller/categories/:id
```

### Authentication

Seller required.

### Request

Category update DTO.

### Success Response

```json
{
  "message": "category edited successfully"
}
```

---

## 14. Delete Category

```http
DELETE /seller/categories/:id
```

### Authentication

Seller required.

### Success Response

```json
{
  "message": "category deleted successfully"
}
```

---

# 📦 Seller Product APIs

Seller authentication is required.

```http
Authorization: Bearer <SELLER_JWT_TOKEN>
```

---

## 15. Create Product

```http
POST /seller/products
```

### Authentication

Seller required.

### Request

The exact `ProductInput` DTO was not available in the current project context.

The product model used by the cart contains at least product information such as:

```text
Product ID
Name
Image URL
Price
Seller/User ID
```

Use the project's product DTO for the exact request schema.

### Success Response

```json
{
  "message": "product created successfully"
}
```

---

## 16. Get Seller Products

```http
GET /seller/products
```

### Authentication

Seller required.

### Success Response

```json
{
  "message": "products fetched successfully",
  "products": []
}
```

---

## 17. Get My Products

```http
GET /seller/my-products
```

### Authentication

Seller required.

### Success Response

```json
{
  "message": "products fetched successfully",
  "products": []
}
```

---

## 18. Update Product

```http
PUT /seller/products/:id
```

or, depending on the configured route:

```http
PATCH /seller/products/:id
```

### Authentication

Seller required.

### Success Response

```json
{
  "message": "product updated successfully"
}
```

---

## 19. Update Stock

The application supports stock updates.

### Success Response

```json
{
  "message": "stock updated successfully"
}
```

The exact route and request DTO should be taken from the catalogue handler currently configured in the project.

---

## 20. Delete Product

```http
DELETE /seller/products/:id
```

### Authentication

Seller required.

### Success Response

```json
{
  "message": "product deleted successfully"
}
```

---

# 🛒 Cart APIs

Base URL:

```text
http://localhost:9000/users
```

Authentication is required.

---

## 21. Add Product To Cart

### Endpoint

```http
POST /users/cart
```

### Authentication

Required

```http
Authorization: Bearer <JWT_TOKEN>
```

### Request

```json
{
  "product_id": 1,
  "qty": 2
}
```

### Success Response

The endpoint returns the updated cart.

```json
{
  "message": "cart created succesfully",
  "data": [
    {
      "product_id": 1,
      "name": "Product Name",
      "qty": 2,
      "price": 100
    }
  ]
}
```

### Cart Behaviour

If the product already exists:

```text
qty >= 1
     ↓
Update quantity
```

If the product exists and:

```text
qty < 1
     ↓
Remove item
```

If the product does not exist:

```text
qty >= 1
     ↓
Create cart item
```

---

# 🛒 Get Cart

## 22. Get Current User Cart

### Endpoint

```http
GET /users/cart
```

### Authentication

Required.

### Request Body

None.

### Success Response

```json
{
  "message": "cart fetched successfully",
  "cart": [
    {
      "id": 1,
      "product_id": 1,
      "name": "Product Name",
      "image_url": "product-image.jpg",
      "seller_id": 2,
      "price": 100,
      "qty": 2
    }
  ],
  "total_amount": 200
}
```

The total is calculated as:

```text
total_amount =
    Σ (product price × quantity)
```

---

# 💳 Payment APIs

The payment system uses **Stripe Checkout**.

Base URL:

```text
http://localhost:9000
```

---

## 23. Create Payment Session

### Endpoint

```http
GET /payment
```

### Authentication

Required.

```http
Authorization: Bearer <JWT_TOKEN>
```

### Request Body

None.

The API gets the user's cart from the authenticated user and calculates the total amount.

### New Payment Response

```json
{
  "message": "payment session created successfully",
  "payment_url": "https://checkout.stripe.com/..."
}
```

The `payment_url` redirects the user to Stripe Checkout.

---

## Existing Payment Session

If the user already has an open Stripe Checkout Session:

```json
{
  "message": "payment session already exists",
  "payment_url": "https://checkout.stripe.com/..."
}
```

The application checks Stripe to determine whether the existing session is still open before reusing its URL.

---

# 💰 Stripe Payment Flow

```text
User
 │
 ▼
GET /payment
 │
 ▼
Authenticate User
 │
 ▼
Check Existing Payment
 │
 ├── Stripe session OPEN
 │       │
 │       ▼
 │   Return existing URL
 │
 └── No active session
         │
         ▼
      Get Cart
         │
         ▼
    Calculate Amount
         │
         ▼
   Generate Order ID
         │
         ▼
 Create Stripe Checkout
         │
         ▼
 Store Payment
         │
         ▼
 Return Checkout URL
```

Stripe metadata:

```json
{
  "order_id": "6629699514",
  "user_id": "1"
}
```

---

# 🧾 Payment Data

A payment contains information such as:

```json
{
  "user_id": 1,
  "amount": 200,
  "payment_id": "cs_test_...",
  "order_id": "6629699514",
  "payment_url": "https://checkout.stripe.com/...",
  "status": "initial"
}
```

Supported payment statuses:

```text
initial
success
failed
pending
```

---

# 👨‍💼 Seller Registration

## 24. Become Seller

### Endpoint

```http
POST /users/become-seller
```

### Authentication

Required.

### Request

The current service uses these seller-input fields:

```json
{
  "first_name": "avc",
  "last_name": "abc",
  "phone": "9876543210",
  "bank_account_number": "XXXXXXXXXXXX",
  "swift_code": "XXXXXXXX",
  "payment_type": "bank"
}
```

> The exact JSON field names should match the project's `dto.SellerInput` definition.

### Success Response

```json
{
  "message": "Seller application submitted",
  "token": "<NEW_SELLER_JWT_TOKEN>"
}
```

The new JWT contains the seller role.

---

# 📦 Order APIs

---

## 25. Create Order

### Endpoint

```http
POST /users/order
```

### Authentication

Required.

### Request Body

No request body is required.

The service obtains the user's current cart and creates an order from the cart items.

### Success Response

```json
{
  "message": "Order created successfully",
  "order": "6629699514"
}
```

The order reference is stored as a **string**.

Example:

```text
6629699514
```

Using a string for an externally exposed order reference makes it easier to later support identifiers such as:

```text
ORD-6629699514
```

without changing the database/API type.

---

# 📋 Get User Orders

## 26. Get Orders

### Endpoint

```http
GET /users/order
```

### Authentication

Required.

### Request Body

None.

### Success Response

```json
{
  "message": "Orders fetched",
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "status": "success",
      "amount": 200,
      "transaction_id": "txn_123",
      "order_ref_number": "6629699514",
      "payment_id": "cs_test_123",
      "items": []
    }
  ]
}
```

---

# 🔎 Get Order By ID

## 27. Get Order

### Endpoint

```http
GET /users/order/:id
```

Example:

```http
GET /users/order/1
```

### Authentication

Required.

### Request Body

None.

### Success Response

```json
{
  "message": "get order by id",
  "order": {
    "id": 1,
    "user_id": 1,
    "status": "success",
    "amount": 200,
    "transaction_id": "txn_123",
    "order_ref_number": "6629699514",
    "payment_id": "cs_test_123",
    "items": [
      {
        "id": 1,
        "order_id": "6629699514",
        "product_id": 1,
        "name": "Product Name",
        "image_url": "product-image.jpg",
        "seller_id": 2,
        "price": 100,
        "qty": 2
      }
    ]
  }
}
```

---

# 👨‍💼 Seller Order APIs

Seller authorization is required.

```http
Authorization: Bearer <SELLER_JWT_TOKEN>
```

---

## 28. Get Seller Orders

### Endpoint

```http
GET /seller/orders
```

### Authentication

Seller required.

### Request Body

None.

### Response

The endpoint returns seller-related order items.

Example structure:

```json
{
  "message": "success",
  "orders": []
}
```

The exact `SellerOrderDetails` response depends on the current DTO implementation.

---

## 29. Get Seller Order Details

### Endpoint

```http
GET /seller/orders/:id
```

Example:

```http
GET /seller/orders/1
```

### Authentication

Seller required.

### Request Body

None.

### Response

```json
{
  "message": "success",
  "order": {}
}
```

The exact fields depend on `dto.SellerOrderDetails`.

---

# 🔑 Authorization Flow

## Normal User

```text
Authorization: Bearer <USER_JWT>
```

Can access:

```text
/users/profile
/users/cart
/users/order
/payment
```

---

## Seller

```text
Authorization: Bearer <SELLER_JWT>
```

Can access seller APIs:

```text
/seller/products
/seller/my-products
/seller/categories
/seller/orders
```

Seller authorization verifies:

```text
JWT valid
    +
user ID valid
    +
role == SELLER
```

---

# 🔢 OTP Generation

EnovaCart uses `crypto/rand` to generate secure random numbers.

OTP range:

```text
100000 - 999999
```

Therefore every generated OTP is a six-digit number.

Example:

```text
583921
```

---

# 🆔 Order Reference Generation

Order references are generated separately from OTPs.

The application uses a larger range for order references and stores the final value as a string.

Example:

```text
6629699514
```

The separation is intentional:

```text
OTP
↓
int
↓
6 digits

Order Reference
↓
string
↓
10-digit numeric identifier currently
```

Using `string` for order references provides flexibility for future formats.

---

# 🗄️ Domain Models

## User

Conceptually contains:

```json
{
  "id": 1,
  "email": "user@example.com",
  "phone": "9876543210",
  "first_name": "Md",
  "last_name": "Ansari",
  "verified": true,
  "user_type": "user"
}
```

---

## Order

```json
{
  "id": 1,
  "user_id": 1,
  "status": "success",
  "amount": 200,
  "transaction_id": "txn_123",
  "order_ref_number": "6629699514",
  "payment_id": "cs_test_123",
  "items": []
}
```

---

## Order Item

```json
{
  "id": 1,
  "order_id": "6629699514",
  "product_id": 1,
  "name": "Product Name",
  "image_url": "product-image.jpg",
  "seller_id": 2,
  "price": 100,
  "qty": 2
}
```

---

## Payment

```json
{
  "id": 1,
  "user_id": 1,
  "amount": 200,
  "transaction_id": "txn_123",
  "customer_id": "cus_123",
  "payment_id": "cs_test_123",
  "order_id": "6629699514",
  "payment_url": "https://checkout.stripe.com/...",
  "status": "initial"
}
```

---

# ⚙️ Environment Variables

Create a `.env` file:

```env
PORT=9000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecommerce

JWT_SECRET=your_jwt_secret

STRIPE_SECRET_KEY=sk_test_xxxxxxxxxxxxx
STRIPE_SUCCESS_URL=http://localhost:3000/success
STRIPE_CANCEL_URL=http://localhost:3000/cancel

SMTP_HOST=your_smtp_host
SMTP_PORT=587
SMTP_USERNAME=your_username
SMTP_PASSWORD=your_password
```

Never commit secrets to GitHub.

Add:

```text
.env
```

to `.gitignore`.

---

# ▶️ Run Locally

## 1. Clone

```bash
git clone <repository-url>
cd ecommerce-app
```

## 2. Install Dependencies

```bash
go mod tidy
```

## 3. Configure Environment

Create:

```text
.env
```

and configure PostgreSQL, JWT, Stripe, and email settings.

## 4. Start PostgreSQL

Using Docker:

```bash
docker compose up -d
```

## 5. Start Application

```bash
go run .
```

Or use Air:

```bash
air
```

Application:

```text
http://localhost:9000
```

---

# 🧪 Testing

Run tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Format code:

```bash
gofmt -w .
```

---

# 🔄 GitHub Actions

The project uses GitHub Actions for CI.

Typical workflow:

```text
Push / Pull Request
        │
        ▼
   Setup Go
        │
        ▼
 Download Dependencies
        │
        ▼
     gofmt
        │
        ▼
      go vet
        │
        ▼
   go test ./...
        │
        ▼
    go build
```

This ensures that code is checked before changes are merged.

---

# 🐳 Docker

Start services:

```bash
docker compose up -d
```

Check containers:

```bash
docker ps
```

Stop containers:

```bash
docker compose down
```

---

# 🔒 Security

EnovaCart implements:

* bcrypt password hashing
* JWT authentication
* Role-based authorization
* Seller-specific authorization
* Secure OTP generation using `crypto/rand`
* OTP expiration
* Environment-based secrets
* Stripe server-side secret key configuration
* User-specific cart access
* User-specific order access

Production deployments should additionally use Stripe webhook signature verification, rate limiting, stronger request validation, structured logging, and proper secret management.

---

# 📡 API Summary

| Method    | Endpoint                 | Auth   | Purpose                     |
| --------- | ------------------------ | ------ | --------------------------- |
| POST      | `/users/register`        | Public | Register                    |
| POST      | `/users/login`           | Public | Login                       |
| GET       | `/users/verify`          | User   | Generate OTP                |
| POST      | `/users/verify`          | User   | Verify OTP                  |
| POST      | `/users/profile`         | User   | Create profile              |
| GET       | `/users/profile`         | User   | Get profile                 |
| PATCH     | `/users/profile`         | User   | Update profile              |
| POST      | `/users/cart`            | User   | Add/update cart             |
| GET       | `/users/cart`            | User   | Get cart                    |
| POST      | `/users/order`           | User   | Create order                |
| GET       | `/users/order`           | User   | Get orders                  |
| GET       | `/users/order/:id`       | User   | Get order                   |
| POST      | `/users/become-seller`   | User   | Become seller               |
| GET       | `/products`              | Public | Get products                |
| GET       | `/products/:id`          | Public | Get product                 |
| GET       | `/seller/categories`     | Seller | Get categories              |
| GET       | `/seller/categories/:id` | Seller | Get category                |
| POST      | `/seller/categories`     | Seller | Create category             |
| PATCH     | `/seller/categories/:id` | Seller | Update category             |
| DELETE    | `/seller/categories/:id` | Seller | Delete category             |
| POST      | `/seller/products`       | Seller | Create product              |
| GET       | `/seller/products`       | Seller | Get seller products         |
| GET       | `/seller/my-products`    | Seller | Get own products            |
| PUT/PATCH | `/seller/products/:id`   | Seller | Update product              |
| DELETE    | `/seller/products/:id`   | Seller | Delete product              |
| GET       | `/seller/orders`         | Seller | Get seller orders           |
| GET       | `/seller/orders/:id`     | Seller | Get seller order            |
| GET       | `/payment`               | User   | Create/reuse Stripe payment |

---

# 🔄 Complete E-Commerce Flow

```text
                    REGISTER
                       │
                       ▼
                     LOGIN
                       │
                       ▼
                  JWT TOKEN
                       │
          ┌────────────┴────────────┐
          │                         │
          ▼                         ▼
     USER PROFILE                PRODUCTS
          │                         │
          └────────────┬────────────┘
                       ▼
                    CART
                       │
                       ▼
                  GET /payment
                       │
                       ▼
               Stripe Checkout
                       │
                       ▼
                   PAYMENT
                       │
                       ▼
                    ORDER
                       │
                       ▼
               ORDER HISTORY
```

Seller flow:

```text
USER
 │
 ▼
BECOME SELLER
 │
 ▼
SELLER JWT
 │
 ├── Categories
 │
 ├── Products
 │
 └── Orders
```

---

# 🚧 Future Improvements

Potential improvements for the project include:

* Stripe webhook payment confirmation
* Payment/order state synchronization
* Inventory reservation
* Stock validation during checkout
* Transactional order creation
* Redis caching
* API rate limiting
* Background workers
* Prometheus metrics
* Structured logging
* Automated integration tests
* Production Docker image
* Cloud deployment
* Improved database constraints
* Centralized error handling

---

# 🎯 Backend Concepts Demonstrated

EnovaCart demonstrates practical Go backend development concepts including:

* REST API development
* Go interfaces
* Dependency injection
* Repository pattern
* Service layer
* Layered architecture
* PostgreSQL
* GORM
* JWT authentication
* Role-based authorization
* Password hashing
* Secure random OTP generation
* Cart management
* Order management
* Payment integration
* Stripe Checkout
* External service integration
* Docker
* Git/GitHub
* GitHub Actions
* CI
* API testing
* Hot reload with Air

---

# 👨‍💻 Author

**Md Shahid Ansari**

Golang / Backend Developer

**Primary Focus**

```text
Go
REST APIs
PostgreSQL
GORM
JWT
Stripe
Docker
GitHub Actions
```

---

# ⭐ Project

If you find EnovaCart useful, consider giving the repository a ⭐ on GitHub.
