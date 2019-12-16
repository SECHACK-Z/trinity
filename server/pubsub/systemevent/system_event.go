package systemevent

type SystemEventType string

const (
	SERVER_START          SystemEventType = "server start"
	WEBUI_START           SystemEventType = "webUI start"
	DIRECTORS_REGISTER    SystemEventType = "director register"
	NEW_SETTINGS_APPLY    SystemEventType = "new settings apply"
	NEW_CONFIG_SAVE       SystemEventType = "new configuration file save"
	HEALTH_CHECK_REGISTER SystemEventType = "new health check register"
	CD_REGISTER           SystemEventType = "new Repository register for CD"
	APPLICATION_START     SystemEventType = "application start"
	KILL_RECEIVED         SystemEventType = "received kill signal"
	KILL_SUCCESS          SystemEventType = "successfully killed application"
	KILL_FAILED           SystemEventType = "failed to kill application"
	BUILD_FAILED          SystemEventType = "build failed"
	WEBHOOK_RECEIVED      SystemEventType = "webhook from github received"
	REPOSITORY_UPDATED    SystemEventType = "respository is updated"
	ERROR                 SystemEventType = "error"
)
