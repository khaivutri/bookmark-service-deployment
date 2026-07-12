package dto

// ShortenURLRes represents the response payload after URL shortening.
type ShortenURLRes struct {
	// Code is the shortened URL code returned to the client.
	Code 		string 	`json:"code" example:"AbCDeFK"`
	// Message provides a status message for the operation.
	Message 	string 	`json:"message" example:"Shorten URL generated successfully!"`
}