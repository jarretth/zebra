package zebra

import (
    "time"
)

const (
    ERR_TIMEOUT ZebraError = iota
)

const (
    WAIT_NO_TIMEOUT = time.Duration(0)
)

const (
    ZXP13_DPI = 300
)
