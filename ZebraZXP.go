package zebra

import (
    "github.com/jarretth/zebrazxp13"
    "runtime"
    "time"
)

func NewZebraZXP(printerName string) *ZebraZXP {
    printer := &ZebraZXP{
        printerName: printerName,
    }
    runtime.SetFinalizer(printer, (*ZebraZXP).CleanUp)
    return printer
}

func (z *ZebraZXP) Name() string {
    return z.printerName
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

func (z *ZebraZXP) DPI() int {
    return ZXP13_DPI
}

func (z *ZebraZXP) EjectCard() uint {
    z.getPrinterHandle()
    ret, _ := zebrazxp13.ZBRPRNEjectCard(zebrazxp13.Handle(z.printerHandle), z.prn_type)
    return uint(ret)
}

func (z *ZebraZXP) WaitIndefinitelyForPrinter() <-chan bool {
    return z.WaitForPrinter(WAIT_NO_TIMEOUT)
}

func (z *ZebraZXP) WaitForPrinter(timeout time.Duration) <-chan bool {
    done := make(chan bool)
    var elapsed time.Duration = 0
    var wait time.Duration = time.Second / 2
    go func() {
        for {
            time.Sleep(wait)
            elapsed += wait
            if z.IsPrinterReady() == 1 {
                done <- true
                break
            }
            if timeout != WAIT_NO_TIMEOUT {
                if elapsed >= timeout {
                    done <- false
                    panic(ERR_TIMEOUT)
                }
                if (elapsed + wait) > timeout {
                    wait = timeout - elapsed
                    if wait <= time.Millisecond {
                        wait = time.Millisecond
                    }
                }
            }
        }
    }()
    return done
}

func (z *ZebraZXP) In(inches float64) int {
    return Inch(inches).Pixels(z.DPI())
}

func (z *ZebraZXP) Mm(millimeters float64) int {
    return Millimeter(millimeters).Pixels(z.DPI())
}

func (z *ZebraZXP) Px(pixels int) int {
    return Pixel(pixels).Pixels(z.DPI())
}

func (z *ZebraZXP) PrintOneSideCard(frontSide GfxCallback) <-chan bool {
    z.getGraphicsHandle()
    gfxContext := newZXPSeries13GraphicsContext(z.graphicsHandle)
    defer gfxContext.CleanUp()
    defer z.freeGraphicsHandle()

    frontSide(gfxContext, z)
    z.printAndClearGraphicsBuffer()

    return z.WaitIndefinitelyForPrinter()
}

func (z *ZebraZXP) PrintTwoSideCard(frontSide GfxCallback, backSide GfxCallback) <-chan bool {
    z.getGraphicsHandle()
    gfxContext := newZXPSeries13GraphicsContext(z.graphicsHandle)
    defer gfxContext.CleanUp()
    defer z.freeGraphicsHandle()

    frontSide(gfxContext, z)
    z.printAndClearGraphicsBuffer()

    backSide(gfxContext, z)
    z.printAndClearGraphicsBuffer()

    return z.WaitIndefinitelyForPrinter()
}

func (z *ZebraZXP) SupportsOneSidedPrinter() (ret bool) {
    defer func() {
        if r := recover(); r != nil {
            ret = false
        }
    }()
    z.getPrinterHandle()
    ret = true
    return ret
}

func (z *ZebraZXP) GetOneSidedPrinter() OneSideCardPrinter {
    return OneSideCardPrinter(z)
}

func (z *ZebraZXP) SupportsTwoSidedPrinter() (ret bool) {
    defer func() {
        if r := recover(); r != nil {
            ret = false
        }
    }()
    z.getPrinterHandle()
    ret = (z.prn_type == zebrazxp13.TYPE_ZXP3_DUAL_SIDE)
    return ret
}

func (z *ZebraZXP) GetTwoSidedPrinter() TwoSideCardPrinter {
    if !z.SupportsTwoSidedPrinter() {
        return nil
    }
    return TwoSideCardPrinter(z)
}

func (z *ZebraZXP) SupportsMagStripeReader() bool {
    return false
}

func (z *ZebraZXP) GetMagStripeReader() TwoSideCardPrinter {
    return nil
}

func (z *ZebraZXP) SupportsMagStripeWriter() bool {
    return false
}

func (z *ZebraZXP) GetMagStripeWriter() MagStripeWriter {
    return nil
}

func (z *ZebraZXP) SupportsMagStripeReaderWriter() bool {
    return false
}

func (z *ZebraZXP) GetMagStripeReaderWriter() MagStripeReaderWriter {
    return nil
}
