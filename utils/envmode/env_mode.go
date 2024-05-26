package envmode

const (
	PROD_MODE EnvMode = iota
	DEV_MODE
	UNIT_TEST_MODE
)

type EnvMode int

var (
	mode EnvMode
)

func GetEnvMode() EnvMode {
	return mode
}

func Is(m EnvMode) bool {
	return mode == m
}

func IsDevMode() bool {
	return Is(DEV_MODE)
}

func IsProdMode() bool {
	return Is(PROD_MODE)
}

func IsUnitTestMode() bool {
	return Is(UNIT_TEST_MODE)
}
