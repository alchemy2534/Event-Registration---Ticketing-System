package services

// In a larger app, the service layer orchestrates complex business logic across multiple repositories.
// For this straightforward REST API, the handlers are directly interacting with repositories for simplicity.
// However, the structure is established here so services can be expanded if logic grows.

type EventService struct {
	// Add event dependencies here if needed in the future
}

func NewEventService() *EventService {
	return &EventService{}
}
