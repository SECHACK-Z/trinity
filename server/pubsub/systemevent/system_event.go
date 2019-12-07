package systemevent

type SystemEventType string

const (
	SERVER_START          SystemEventType = "server start"
	WEBUI_START           SystemEventType = "webUI start"
	DIRECTORS_REGISTER    SystemEventType = "directors register"
	NEW_SETTINGS_APPLY    SystemEventType = "new settings apply"
	NEW_CONFIG_SAVE       SystemEventType = "new configuration file save"
	HEALTH_CHECK_REGISTER SystemEventType = "new health check register"
	CD_REGISTER           SystemEventType = "new Repository register for CD"
	ERROR                 SystemEventType = "error"
)
