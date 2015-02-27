package zebra
import (
    "github.com/jarretth/zebrazxp13"
    "runtime"
    "time"
)

func NewZebraZXP(printerName string) *ZebraZXP {
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

func (z *ZebraZXP) DPI() int {
    return ZXP13_DPI
}

func (z *ZebraZXP) EjectCard() uint {
    z.getPrinterHandle()
    ret, _ := zebrazxp13.ZBRPRNEjectCard(zebrazxp13.Handle(z.printerHandle), z.prn_type)
    return uint(ret)
}

func (z *ZebraZXP) WaitIndefinitelyForPrinter() {
    z.WaitForPrinter(WAIT_NO_TIMEOUT)
}

func (z *ZebraZXP) WaitForPrinter(timeout time.Duration) {
    var elapsed time.Duration = 0
    var wait time.Duration    = time.Second / 2
    for ;; {
        time.Sleep(wait)
        elapsed += wait
        if z.IsPrinterReady() == 1 {
            break
        }
        if timeout != WAIT_NO_TIMEOUT {
            if elapsed >= timeout {
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
}

func (z *ZebraZXP) PrintOneSideCard(frontSide GfxCallback) uint {
    z.getGraphicsHandle()
    gfxContext := newZXPSeries13GraphicsContext(z.graphicsHandle)
    defer gfxContext.CleanUp()
    defer z.freeGraphicsHandle()

    frontSide(gfxContext)
    z.printAndClearGraphicsBuffer()

    return 1
}

func (z *ZebraZXP) PrintTwoSideCard(frontSide GfxCallback, backSide GfxCallback) uint {
    z.getGraphicsHandle()
    gfxContext := newZXPSeries13GraphicsContext(z.graphicsHandle)
    defer gfxContext.CleanUp()
    defer z.freeGraphicsHandle()

    frontSide(gfxContext)
    z.printAndClearGraphicsBuffer()

    backSide(gfxContext)
    z.printAndClearGraphicsBuffer()

    return 1
}

func (z *ZebraZXP) SupportsOneSidedPrinter() bool {
    return true
}

func (z *ZebraZXP) GetOneSidedPrinter() OneSideCardPrinter {
    return OneSideCardPrinter(z)
}

func (z *ZebraZXP) SupportsTwoSidedPrinter() bool {
    z.getPrinterHandle()
    return z.prn_type == zebrazxp13.TYPE_ZXP3_DUAL_SIDE
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
