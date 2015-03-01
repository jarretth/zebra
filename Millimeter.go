package zebra

import "math"

func (m Millimeter) Pixels(DPI int) int {
    return int(math.Ceil((float64(m) / 25.4) * float64(DPI)))
}
