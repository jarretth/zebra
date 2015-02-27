package zebra

import (
    "github.com/jarretth/zebrazxp13"
    "image"
    "image/color"
    "time"
)

type PrinterHandle zebrazxp13.Handle
type GraphicsHandle zebrazxp13.Handle

type ZebraZXP struct {
    printerName string
    prn_type zebrazxp13.PrinterType
    printerHandle PrinterHandle
    graphicsHandle GraphicsHandle
}

type ZebraZXPSeries13GfxContext struct {
    graphicsHandle GraphicsHandle
}

type GfxCallback func(GfxContext)

type GfxContext interface {
    DrawBarcode(location image.Point, rotation zebrazxp13.BarCodeRotation, barcodetype zebrazxp13.BarCodeType, barwidthratio zebrazxp13.BarCodeWidth, barcodemultiplier int, barcodeheight int, textunder zebrazxp13.BarCodeTextUnder, barcodedata string) (ret int)
    DrawText(x uint, y uint, text string, font string, fontsize uint, fontstyle zebrazxp13.TextFontStyle, color color.Color) (ret int)
    DrawTextRectangle(rect image.Rectangle, alignment zebrazxp13.TextAlignment, text string, font string, fontsize uint, fontstyle zebrazxp13.TextFontStyle, color color.Color) (ret int)
    DrawLine(start image.Point, end image.Point, color color.Color, thickness float32) (ret int)
    DrawRectangle(rect image.Rectangle, color color.Color, thickness float32) (ret int)
    DrawEllipse(rect image.Rectangle, color color.Color, thickness float32) (ret int)
}

type Measurement interface {
    Pixels(DPI int) int
}

type Pixel struct {
    dot int
}

type Inch struct {
    dot float64
}

type Millimeter struct {
    dot float64
}

type Printer interface {
    CleanUp()

    DPI() int

    IsPrinterReady() uint
    WaitForPrinter(timeout time.Duration)
    WaitIndefinitelyForPrinter()

    SupportsOneSidedPrinter() bool
    GetOneSidedPrinter() OneSideCardPrinter

    SupportsTwoSidedPrinter() bool
    GetTwoSidedPrinter() TwoSideCardPrinter

    SupportsMagStripeReader() bool
    GetMagStripeReader() TwoSideCardPrinter

    SupportsMagStripeWriter() bool
    GetMagStripeWriter() MagStripeWriter

    SupportsMagStripeReaderWriter() bool
    GetMagStripeReaderWriter() MagStripeReaderWriter
}

type OneSideCardPrinter interface {
    Printer
    EjectCard() uint
    PrintOneSideCard(frontSide GfxCallback) uint
}

type TwoSideCardPrinter interface {
    OneSideCardPrinter
    PrintTwoSideCard(frontSide GfxCallback, backSide GfxCallback) uint
}

type MagStripeReader interface {

}

type MagStripeWriter interface {

}

type MagStripeReaderWriter interface {
    MagStripeReader
    MagStripeWriter
}
