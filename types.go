package zebra

import (
    "github.com/jarretth/zebrazxp13"
)

type PrinterHandle zebrazxp13.Handle
type GraphicsHandle zebrazxp13.Handle
type Pixel int
type Inch float64
type Millimeter float64
type GfxCallback func(GfxContext, Printer)
type ZebraError int

type ZebraZXP struct {
    printerName    string
    prn_type       zebrazxp13.PrinterType
    printerHandle  PrinterHandle
    graphicsHandle GraphicsHandle
}

type ZebraZXPSeries13GfxContext struct {
    graphicsHandle GraphicsHandle
}
