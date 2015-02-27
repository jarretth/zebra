package zebra

import (
    "fmt"
    "github.com/jarretth/zebrazxp13"
    "image/color"
    "image"
)

func getColor(color color.Color) uint32 {
    r,  g, b, _ := color.RGBA()
    return uint32((r << 16) | (g << 8) | b);
}

func GetPrinterSDKVersion() (version string) {
    major,minor,engLevel := zebrazxp13.ZBRPRNGetSDKVer()
    return fmt.Sprintf("%d.%d.%d", major, minor, engLevel)
}

func GetGraphicsSDKVersion() (version string) {
    major,minor,engLevel := zebrazxp13.ZBRGDIGetSDKVer()
    return fmt.Sprintf("%d.%d.%d", major, minor, engLevel)
}

func Rect(x1 Measurement, y1 Measurement, x2 Measurement, y2 Measurement, p Printer) image.Rectangle {
    DPI := p.DPI()
    return image.Rect(x1.Pixels(DPI), y1.Pixels(DPI), x2.Pixels(DPI), y2.Pixels(DPI))
}

func Pt(x Measurement, y Measurement, p Printer) image.Point {
    DPI := p.DPI()
    return image.Pt(x.Pixels(DPI), y.Pixels(DPI))
}
