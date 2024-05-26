package internal

import (
	"fmt"
	"os"

	"tradeengine/utils/logger/color"

	"go.uber.org/zap/zapcore"
)

func newColorConsoleCore(l *MyLogger) zapcore.Core {
	levelEncoder := capitalColorLevelEncoderWithConfig(l)
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",

		EncodeTime:   zapcore.TimeEncoderOfLayout("2006/01/02-15:04:05.000"),
		EncodeLevel:  levelEncoder,
		EncodeCaller: goroutineCallerEncoder,
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		l.atomicLevel,
	)
}

func capitalColorLevelEncoderWithConfig(lg *MyLogger) func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	isLogColorCustom := lg.property.isLogColorCustom
	if !isLogColorCustom {
		return capitalColorLevelEncoder
	}
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		ct := color.New(fmt.Sprintf("[%-5s]", l.CapitalString()))
		var levelKey string
		for _, isFg := range []bool{true, false} {
			switch l {
			case zapcore.DebugLevel:
				levelKey = "debug"
			case zapcore.InfoLevel:
				levelKey = "info"
			case zapcore.WarnLevel:
				levelKey = "warning"
			case zapcore.ErrorLevel:
				levelKey = "error"
			case zapcore.DPanicLevel:
				levelKey = "dpanic"
			case zapcore.PanicLevel:
				levelKey = "panic"
			case zapcore.FatalLevel:
				levelKey = "fatal"
			default:
				levelKey = ""
			}
			var pos string
			if isFg {
				pos = "fg"
			} else {
				pos = "bg"
			}
			code, ok := lg.property.logColor[levelKey+"_"+pos]
			if !ok {
				code = 0x000000FF
			}
			isCustom, isSGR, colors := logColorParse(code)
			if isCustom {
				applyLogColor(ct, isFg, isSGR, colors)
			} else {
				switch l {
				case zapcore.DebugLevel:
					if isFg {
						ct.SetFgObj(color.GetDefaultFgDebug())
					} else {
						ct.SetBgObj(color.GetDefaultBgDebug())
					}
				case zapcore.InfoLevel:
					if isFg {
						ct.SetFgObj(color.GetDefaultFgInfo())
					} else {
						ct.SetBgObj(color.GetDefaultBgInfo())
					}
				case zapcore.WarnLevel:
					if isFg {
						ct.SetFgObj(color.GetDefaultFgWarning())
					} else {
						ct.SetBgObj(color.GetDefaultBgWarning())
					}
				case zapcore.ErrorLevel:
					if isFg {
						ct.SetFgObj(color.GetDefaultFgError())
					} else {
						ct.SetBgObj(color.GetDefaultBgError())
					}
				case zapcore.DPanicLevel:
					if isFg {
						ct.SetSfg(color.FG_Yellow)
					} else {
						ct.SetSbg(color.BG_Red)
					}
				case zapcore.PanicLevel:
					if isFg {
						ct.SetSfg(color.FG_Yellow)
					} else {
						ct.SetSbg(color.BG_Red)
					}
				case zapcore.FatalLevel:
					if isFg {
						ct.SetSfg(color.FG_Yellow)
					} else {
						ct.SetSbg(color.BG_Red)
					}
				default:
					if isFg {
						ct.SetSfg(color.FG_Red)
					}
				}
			}
		}
		enc.AppendString(ct.Build())
	}
}

func logColorParse(code int) (isCustom bool, isSGR bool, colors []uint8) {
	if code&0xFF == 0xFF {
		return false, false, []uint8{}
	}
	if code&0xFFFFFF00 == 0 && code&0xFF > 0 {
		return true, true, []uint8{uint8(code & 0xFF)}
	}
	return true, false, []uint8{
		uint8(code & 0xFF000000),
		uint8(code & 0x00FF0000),
		uint8(code & 0x0000FF00),
	}
}

func applyLogColor(ct *color.ColorText, isFg bool, isSGR bool, colors []uint8) {
	if isSGR {
		if isFg {
			ct.SetSfg(color.FgColor(colors[0]))
		} else {
			ct.SetSbg(color.BgColor(colors[0]))
		}
	} else {
		if isFg {
			ct.SetFg(colors[0], colors[1], colors[2])
		} else {
			ct.SetBg(colors[0], colors[1], colors[2])
		}
	}
}

func capitalColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	ct := color.New(fmt.Sprintf("[%-5s]", l.CapitalString()))
	switch l {
	case zapcore.DebugLevel:
		ct.SetFgObj(color.GetDefaultFgDebug()).SetBgObj(color.GetDefaultBgDebug())
	case zapcore.InfoLevel:
		ct.SetFgObj(color.GetDefaultFgInfo()).SetBgObj(color.GetDefaultBgInfo())
	case zapcore.WarnLevel:
		ct.SetFgObj(color.GetDefaultFgWarning()).SetBgObj(color.GetDefaultBgWarning())
	case zapcore.ErrorLevel:
		ct.SetFgObj(color.GetDefaultFgError()).SetBgObj(color.GetDefaultBgError())
	case zapcore.DPanicLevel:
		ct.SetSfg(color.FG_Yellow).SetSbg(color.BG_Red)
	case zapcore.PanicLevel:
		ct.SetSfg(color.FG_Yellow).SetSbg(color.BG_Red)
	case zapcore.FatalLevel:
		ct.SetSfg(color.FG_Yellow).SetSbg(color.BG_Red)
	default:
		ct.SetSfg(color.FG_Red)
	}
	enc.AppendString(ct.Build())
}

func goroutineCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	goroutineText := fmt.Sprintf("[GO-%d]", getGoID())
	enc.AppendString(goroutineText)
	text := fmt.Sprintf("[%s]", caller.TrimmedPath())
	enc.AppendString(text)
}
