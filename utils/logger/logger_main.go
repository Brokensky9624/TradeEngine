package logger

import "tradeengine/utils/logger/internal"

/*
Usage:

	These logger can be directly used, no need to initialize.

	Ex: logger.REST.Debug("..%d..%s", int, string)  // for trival information (TRACE)
	Ex: logger.REST.Info("..%d..%s", int, string)   // for debug content & important inforamtion (DEBUG + INFO)
	Ex: logger.REST.Warn("..%d..%s", int, string)   // for warning information
	Ex: logger.REST.Error("..%d..%s", int, string)  // for error message
*/

var SERVER = internal.NewMyLogger("server_log.properties")
var REST = internal.NewMyLogger("rest_log.properties")
var DB = internal.NewMyLogger("db_log.properties")
var Panic = internal.NewMyLogger("panic_log.properties")

/*
only write to stdout
*/
var STD = internal.NewMyLogger("std_log.properties")
