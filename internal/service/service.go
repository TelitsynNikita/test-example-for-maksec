package service

type AgentConfig struct {
	CallbackURL   string
	CallbackToken string
}

type Services struct {
	ScriptService *ScriptService
	EventService  *EventService
}
