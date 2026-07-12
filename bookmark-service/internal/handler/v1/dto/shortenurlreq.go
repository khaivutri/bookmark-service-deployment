package dto

// ShortenURLReq represents the request payload for URL shortening.
type ShortenURLReq struct {
	// URL is the original URL to shorten.
	URL 	string 	`json:"url" binding:"url,required" example:"https://example.com"`
	// Exp is the expiration time of the shortened URL in seconds.
	Exp 	int64 	`json:"exp" binding:"required,gte=5"`
}