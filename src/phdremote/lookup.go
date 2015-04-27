package phdremote

import "fmt"

func LookupObject(object string) (string, string) {
    messier, ok := Messier[object]
    if (ok)  {
        return fmt.Sprintf("%f", messier[0]), fmt.Sprintf("%f", messier[1])
    }
    ngc, ok := NGC[object]
    if (ok)  {
        return fmt.Sprintf("%f", ngc[0]), fmt.Sprintf("%f", ngc[1])
    }
    caldwell, ok := Caldwell[object]
    if (ok)  {
        ngc, ok = NGC[caldwell]
        return fmt.Sprintf("%f", ngc[0]), fmt.Sprintf("%f", ngc[1])
    }
    star, ok := Stars[object]
    if (ok)  {
        return fmt.Sprintf("%f", star[0]), fmt.Sprintf("%f", star[1])
    }
    return "", ""
}
