package zebra

import "math"

func (i Inch) Pixels(DPI int) int {
    return int(math.Ceil(float64(i) * float64(DPI)))
}
