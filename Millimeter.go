package zebra

import "math"

func Mm(dot float64) (m *Millimeter) {
    m = &Millimeter{
        dot: dot,
    }
    return m
}

func (m *Millimeter) Pixels(DPI int) int {
    return int(math.Ceil((m.dot / 25.4) * float64(DPI)))
}
