package derr

var (
	BadRequestError     = NewClientError("BAD_REQUEST", "Bad Request")
	UnauthorizedError   = NewClientError("UNAUTHORIZED", "Unauthorized")
	NotFoundError       = NewClientError("NOT FOUND", "Not Found")
	InternalServerError = NewRepositoryError("INTERNAL_SERVER_ERROR", "Internal Server Error")
)

var (
	InvalidParameterError = NewClientError("INVALID_PARAMETER", "Invalid Parameter")
)

var (
	InvalidCredentials = NewClientError("INVALID_CREDENTIALS", "Invalid credentials")
	UserAlreadyExists  = NewClientError("USER_ALREADY_EXISTS", "User already exists")
)

var (
	InvalidProductName        = NewClientError("INVALID_PRODUCT_NAME", "Invalid product name")
	InvalidProductDescription = NewClientError("INVALID_DESCRIPTION", "Invalid product description")
	InvalidProductPrice       = NewClientError("INVALID_PRODUCT_PRICE", "Invalid product price")
	InvalidProductStock       = NewClientError("INVALID_PRODUCT_STOCK", "Invalid product stock")
	InvalidProductDiscount    = NewClientError("INVALID_PRODUCT_DISCOUNT", "Invalid product discount")
)
