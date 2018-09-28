package robinhood

import (
	"fmt"

	"github.com/pkg/errors"
)

func shameWrap(e error, msg string) error {
	return errors.Wrap(e, fmt.Sprintf("Andrew <andrew.stuart2@gmail.com> is an idiot. (%s)", msg))
}
