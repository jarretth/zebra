package zebra

import (
    "github.com/jarretth/zebrazxp13"
    "image"
    "image/color"
    "unicode/utf16"
)

func convertToUtf16(i string) string {
    b := utf16.Encode([]rune(i))
    output := make([]byte, 0)
    for i := range b {
        var h, l byte = byte(b[i] >> 8), byte(b[i] & 0xFF)
        if h == 0 {
            output = append(output, l)
        } else {
            output = append(output, h)
            output = append(output, l)
        }
    }
    return string(append(output, []byte{0, 0}...))
}

func newZXPSeries13GraphicsContext(handle GraphicsHandle) *ZebraZXPSeries13GfxContext {
    context := &ZebraZXPSeries13GfxContext{
        graphicsHandle: handle,
    }
    return context
}

func (g *ZebraZXPSeries13GfxContext) DrawBarcode(location image.Point, rotation zebrazxp13.BarCodeRotation, barcodetype zebrazxp13.BarCodeType, barwidthratio zebrazxp13.BarCodeWidth, barcodemultiplier int, barcodeheight int, textunder zebrazxp13.BarCodeTextUnder, barcodedata string) (ret int) {
    uret, err := zebrazxp13.ZBRGDIDrawBarCode(
        uint(location.X),
        uint(location.Y),
        rotation,
        barcodetype,
        barwidthratio,
        uint(barcodemultiplier),
        uint(barcodeheight),
        textunder,
        barcodedata,
    )
    if uret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    return int(uret)
}

func (g *ZebraZXPSeries13GfxContext) DrawText(x uint, y uint, text string, font string, fontsize uint, fontstyle zebrazxp13.TextFontStyle, color color.Color) (ret int) {
    // uret, err := zebrazxp13.ZBRGDIDrawText(
    uret, err := zebrazxp13.ZBRGDIDrawTextUnicode(
        x,
        y,
        text,
        font,
        fontsize,
        fontstyle,
        getColor(color),
    )
    if uret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    return int(uret)
}

func (g *ZebraZXPSeries13GfxContext) DrawTextRectangle(rect image.Rectangle, alignment zebrazxp13.TextAlignment, text string, font string, fontsize uint, fontstyle zebrazxp13.TextFontStyle, color color.Color) (ret int) {
    uret, err := zebrazxp13.ZBRGDIDrawTextRect(
        uint(rect.Min.X),
        uint(rect.Min.Y),
        uint(rect.Dx()),
        uint(rect.Dy()),
        alignment,
        convertToUtf16(text),
        font,
        fontsize,
        fontstyle,
        getColor(color),
    )
    if uret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    return int(uret)
}

func (g *ZebraZXPSeries13GfxContext) DrawLine(start image.Point, end image.Point, color color.Color, thickness float32) (ret int) {
    uret, err := zebrazxp13.ZBRGDIDrawLine(
        uint(start.X),
        uint(start.Y),
        uint(end.X),
        uint(end.Y),
        getColor(color),
        thickness,
    )
    if uret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    return int(uret)
}

func (g *ZebraZXPSeries13GfxContext) DrawRectangle(rect image.Rectangle, color color.Color, thickness float32) (ret int) {
    uret, err := zebrazxp13.ZBRGDIDrawRectangle(
        uint(rect.Min.X),
        uint(rect.Min.Y),
        uint(rect.Dx()),
        uint(rect.Dy()),
        thickness,
        getColor(color),
    )
    if uret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    return int(uret)
}

func (g *ZebraZXPSeries13GfxContext) DrawEllipse(rect image.Rectangle, color color.Color, thickness float32) (ret int) {
    uret, err := zebrazxp13.ZBRGDIDrawEllipse(
        uint(rect.Min.X),
        uint(rect.Min.Y),
        uint(rect.Dx()),
        uint(rect.Dy()),
        thickness,
        getColor(color),
    )
    if uret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    return int(uret)
}

func (g *ZebraZXPSeries13GfxContext) CleanUp() {

}
