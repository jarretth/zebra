package zebra

import "math"

func In(dot float64) (i *Inch) {
    i = &Inch{
        dot: dot,
    }
    return i
}

func (i *Inch) Pixels(DPI int) int {
    return int(math.Ceil(i.dot * float64(DPI)))
}
