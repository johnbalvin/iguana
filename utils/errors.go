package utils

import "errors"

var ErrDepth = errors.New("Deph over the maximum allow: 30, this would mean you have a circular depencie, or have more than 30 nested dependencies")
var ErrUnknown = errors.New("The obfuscator throwed an unknown error")
