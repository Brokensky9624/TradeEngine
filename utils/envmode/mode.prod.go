//go:build prod
// +build prod

package envmode

func init() {
	mode = PROD_MODE
}
