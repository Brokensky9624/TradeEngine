//go:build !prod
// +build !prod

package envmode

func init() {
	mode = DEV_MODE
}
