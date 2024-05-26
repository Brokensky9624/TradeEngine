package shared

type IMylogger interface {
	Debug(template string, a ...interface{})
	Info(template string, a ...interface{})
	Warn(template string, a ...interface{})
	Error(template string, a ...interface{})
}
