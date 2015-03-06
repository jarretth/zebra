package zebra

func (z ZebraError) Error() string {
    switch z {
    case ERR_TIMEOUT:
        return "Timed out waiting for printer"
    }
    return "An error occurred"
}
