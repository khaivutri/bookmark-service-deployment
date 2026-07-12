package model

// HealthReport represents the health status response.
type HealthReport struct {
	Message 		string 				`json:"message" example:"OK"`
	ServiceName 	string 				`json:"service_name" example:"bookmark_service"`
	InstanceID 		string				`json:"instance_id" example:"cbe1a562-596b-45d0-bf8b-a999b23b184a"`
	Dependencies	map[string]string	`json:"dependency" example:"{\"redis\": \"UP\"}"`
}