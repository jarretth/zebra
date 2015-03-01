package zebra

import (
    "fmt"
    "github.com/jarretth/zebrazxp13"
    "image/color"
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
