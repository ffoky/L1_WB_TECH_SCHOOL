package cut

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errNoFields = errors.New("no fields specified")

// разворачивает строки 1,3-5 в 1, 3, 4, 5
func ParseFields(fields string) ([]int, error) {
	if fields == "" {
		return nil, errNoFields
	}

	parts := strings.Split(fields, ",")
	res := make([]int, 0)
	for _, p := range parts {
		if strings.Contains(p, "-") {
			rangeParts := strings.Split(p, "-")
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("bad field: %s", p)
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("bad field: %s", p)
			}
			for i := start; i <= end; i++ {
				res = append(res, i)
			}
			continue
		}
		num, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("bad field: %s", p)
		}
		res = append(res, num)
	}
	return res, nil
}
