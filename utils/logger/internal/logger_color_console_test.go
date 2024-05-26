package internal_test

import (
	"log"
	"testing"

	"tradeengine/utils/logger"
	"tradeengine/utils/logger/color"
)

func Test_fgColor_AddWithBg(t *testing.T) {
	logger.REST.Debug("DebugLevel")
	logger.REST.Info("InfoLevel")
	logger.REST.Warn("WarnLevel")
	logger.REST.Error("ErrorLevel")

	log.Println(color.New("DPanicLevel").SetSfg(color.FG_Yellow).SetSbg(color.BG_Red).Build())
	log.Println(color.New("PanicLevel").SetSfg(color.FG_Yellow).SetSbg(color.BG_Red).Build())
	log.Println(color.New("FatalLevel").SetSfg(color.FG_Yellow).SetSbg(color.BG_Red).Build())

	log.Println(color.New("Default_Debug").SetBgObj(color.GetDefaultBgDebug()).SetFgObj(color.GetDefaultFgDebug()).Build())
	log.Println(color.New("Default_Info").SetBgObj(color.GetDefaultBgInfo()).SetFgObj(color.GetDefaultFgInfo()).Build())
	log.Println(color.New("Default_Warning").SetBgObj(color.GetDefaultBgWarning()).SetFgObj(color.GetDefaultFgWarning()).Build())
	log.Println(color.New("Default_Error").SetBgObj(color.GetDefaultBgError()).SetFgObj(color.GetDefaultFgError()).Build())
}
