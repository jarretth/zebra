package zebra

func Px(dot int) (p *Pixel) {
    p = &Pixel{
        dot: dot,
    }
    return p
}

func (p *Pixel) Pixels(DPI int) int {
    return p.dot
}
