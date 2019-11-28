package systemevent

type SystemEventType string

const (
	SERVER_START       SystemEventType = "server start"
	WEBUI_START        SystemEventType = "webUI start"
	DIRECTORS_REGISTER SystemEventType = "directors register"
	NEW_SETTINGS_APPLY SystemEventType = "new settings apply"
	NEW_CONFIG_SAVE    SystemEventType = "new configuration file save"

	ERROR SystemEventType = "error"
)
