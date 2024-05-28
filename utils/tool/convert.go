package tool

import "strconv"

func ParseUintFromStr(idStr string) (uint, error) {
	n, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
