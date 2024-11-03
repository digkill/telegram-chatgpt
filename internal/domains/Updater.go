package domains

type UpdateTelegramData struct {
	*Handler
}

type Updater interface{}

func NewUpdateTelegramData(handler *Handler) *UpdateTelegramData {
	return &UpdateTelegramData{handler}
}
