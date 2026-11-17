package messages

const (
	ErrorInvalidRequest = "invalid_request"
	ErrorUnauthorized   = "unauthorized"
	ErrorNotFound       = "not_found"
	ErrorForbidden      = "forbidden"
	ErrorAlreadyExists  = "already_exists"
	ErrorLimitReached   = "limit_reached"

	AuthorizationHeaderRequired = "Authorization header required"
	InvalidAuthorization        = "Invalid or expired token"
	AdvertisementNotFound       = "Advertisement not found"
	FavoriteNotFound            = "Favorite entry not found"
	InvalidAdvertisementId      = "Invalid advertisement ID"
	AlreadyInFavorites          = "Advertisement is already in favorites"
	CannotAddOwnAdvertisement   = "Cannot add own advertisement to favorites"
	FavoritesLimitReached       = "Favorites limit of 50 reached"
	AdvertisementIdRequired     = "Key: 'AddFavoriteRequest.AdvertisementID' Error:Field validation for 'AdvertisementID' failed on the 'required' tag"
)
