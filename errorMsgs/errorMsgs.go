package errorMsgs

import "errors"

var DictNotFound = errors.New("Not Found word.")

var DictAlreadyRegistered = errors.New("Already Registered.")
