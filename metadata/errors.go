package metadata

import "errors"

var NothingToRead = errors.New("Nothing to read")
var InvalidPacketError = errors.New("Invalid Packet")
var TimedOut = errors.New("Timed out while trying to read metadata")
