package derr

var (
	BadRequestError     = NewBadRequestError("Bad Request")
	UnauthorizedError   = NewUnauthorizedError("Unauthorized")
	NotFoundError       = NewNotFoundError("Not Found")
	InternalServerError = NewInternalError("Internal Server Error")
)

var (
	InvalidCredentials = NewClientError("INVALID_CREDENTIALS", "Invalid credentials")
	InvalidEmail       = NewClientError("INVALID_EMAIL", "Invalid email")
	UserAlreadyExists  = NewClientError("USER_ALREADY_EXISTS", "User already exists")
)

var (
	InvalidUsername    = NewClientError("INVALID_USERNAME", "Invalid username")
	BiographyIsTooLong = NewClientError("BIOGRAPHY_IS_TOO_LONG", "Biography is too long")
	NameIsTooShort     = NewClientError("NAME_IS_TOO_SHORT", "Name is too short")
	NameIsTooLong      = NewClientError("NAME_IS_TOO_LONG", "Name is too long")
	WeakPassword       = NewClientError("WEAK_PASSWORD", "Weak password")
)

var (
	InvalidProductName        = NewClientError("INVALID_PRODUCT_NAME", "Invalid product name")
	InvalidProductDescription = NewClientError("INVALID_DESCRIPTION", "Invalid product description")
	InvalidProductPrice       = NewClientError("INVALID_PRODUCT_PRICE", "Invalid product price")
	InvalidProductStock       = NewClientError("INVALID_PRODUCT_STOCK", "Invalid product stock")
	InvalidProductDiscount    = NewClientError("INVALID_PRODUCT_DISCOUNT", "Invalid product discount")
)
var (
	InvalidActivityType     = NewClientError("INVALID_ACTIVITY_TYPE", "Invalid activity type")
	InvalidActivityDuration = NewClientError("INVALID_ACTIVITY_DURATION", "Invalid activity duration")
)
