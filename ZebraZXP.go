package zebra
import (
    "fmt"
    "image"
    "image/color"
    "runtime"
    "time"
    "github.com/jarretth/zebrazxp13"
)

type PrinterHandle zebrazxp13.Handle
type GraphicsHandle zebrazxp13.Handle

type ZebraZXP struct {
    printerName string
    prn_type uint
    printerHandle PrinterHandle
    graphicsHandle GraphicsHandle
}

type GfxContext interface {
    DrawBarcode(location image.Point, rotation zebrazxp13.BarCodeRotation, barcodetype zebrazxp13.BarCodeType, barwidthratio zebrazxp13.BarCodeWidth, barcodemultiplier int, barcodeheight int, textunder zebrazxp13.BarCodeTextUnder, barcodedata string) (ret int)
    DrawText(x uint, y uint, text string, font string, fontsize uint, fontstyle zebrazxp13.TextFontStyle, color color.Color) (ret int)
    DrawTextRectangle(rect image.Rectangle, alignment zebrazxp13.TextAlignment, text string, font string, fontsize uint, fontstyle zebrazxp13.TextFontStyle, color color.Color) (ret int)
    DrawLine(start image.Point, end image.Point, color color.Color, thickness float32) (ret int)
    DrawRectangle(rect image.Rectangle, color color.Color, thickness float32) (ret int)
    DrawEllipse(rect image.Rectangle, color color.Color, thickness float32) (ret int)
}

type GfxCallback func(GfxContext)

func New(printerName string) *ZebraZXP {
    printer := &ZebraZXP {
        printerName: printerName,
    }
    runtime.SetFinalizer(printer, (*ZebraZXP).CleanUp)
    return printer
}

func (z *ZebraZXP) getGraphicsHandle() {
    if z.graphicsHandle != 0 {
        return
    }
    ret, graphicsHandle, err := zebrazxp13.ZBRGDIInitGraphics(z.printerName)
    if ret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    z.graphicsHandle = GraphicsHandle(graphicsHandle)
}

func (z *ZebraZXP) freeGraphicsHandle() {
    if z.graphicsHandle == 0 {
        return
    }
    defer func() {
        z.graphicsHandle = 0
    }()
    ret, err := zebrazxp13.ZBRGDICloseGraphics(zebrazxp13.Handle(z.graphicsHandle))
    if ret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
}

func (z *ZebraZXP) getPrinterHandle() {
    if z.printerHandle != 0 {
        return
    }
    ret, printerHandle, prn_type, err := zebrazxp13.ZBRGetHandle(z.printerName)
    if ret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    z.printerHandle = PrinterHandle(printerHandle)
    z.prn_type = prn_type
}

func (z *ZebraZXP) freePrinterHandle() {
    if z.printerHandle == 0 {
        return
    }

    defer func() {
        z.printerHandle = 0
        z.prn_type = 0
    }()

    ret, err := zebrazxp13.ZBRCloseHandle(zebrazxp13.Handle(z.printerHandle))
    if ret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
}

func (z *ZebraZXP) printAndClearGraphicsBuffer() {
    ret, err := zebrazxp13.ZBRGDIPrintGraphics(zebrazxp13.Handle(z.graphicsHandle))
    if ret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
    ret, err = zebrazxp13.ZBRGDIClearGraphics()
    if ret != zebrazxp13.ZBR_SUCCESS {
        panic(err)
    }
}

func (z *ZebraZXP) CleanUp() {
    defer z.freePrinterHandle()
    defer z.freeGraphicsHandle()
}

func (z *ZebraZXP) IsPrinterReady() uint {
    ret, err := zebrazxp13.ZBRGDIIsPrinterReady(z.printerName)
    if err != 0 {
        panic(err)
    }
    return uint(ret)
}

func (z *ZebraZXP) EjectCard() uint {
    z.getPrinterHandle()
    ret, _ := zebrazxp13.ZBRPRNEjectCard(zebrazxp13.Handle(z.printerHandle), z.prn_type)
    return uint(ret)
}

func (z *ZebraZXP) WaitForPrinter() {
    for ;; {
        time.Sleep(time.Second/2)
        if z.IsPrinterReady() == 1 {
            break
        }
    }
}

func (z *ZebraZXP) PrintOneSideCard(frontSide GfxCallback) uint {
    z.getGraphicsHandle()
    gfxContext := newZXPSeries13GraphicsContext(z.graphicsHandle)
    defer gfxContext.CleanUp()

    frontSide(gfxContext)
    z.printAndClearGraphicsBuffer()

    return 1
}

func (z *ZebraZXP) PrintTwoSideCard(frontSide GfxCallback, backSide GfxCallback) uint {
    z.getGraphicsHandle()
    gfxContext := newZXPSeries13GraphicsContext(z.graphicsHandle)
    defer gfxContext.CleanUp()

    frontSide(gfxContext)
    z.printAndClearGraphicsBuffer()

    backSide(gfxContext)
    z.printAndClearGraphicsBuffer()

    return 1
}

func GetPrinterSDKVersion() (version string) {
    major,minor,engLevel := zebrazxp13.ZBRPRNGetSDKVer()
    return fmt.Sprintf("%d.%d.%d", major, minor, engLevel)
}

func GetGraphicsSDKVersion() (version string) {
    major,minor,engLevel := zebrazxp13.ZBRGDIGetSDKVer()
    return fmt.Sprintf("%d.%d.%d", major, minor, engLevel)
}
