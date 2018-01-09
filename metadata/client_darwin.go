// +build darwin

package metadata

type SerialClient struct {
}

func (client *SerialClient) Connect(serialPort string) error {
	panic("not supported on the current operating system")
}

func (client *SerialClient) WriteLine(line string) error {
	panic("not supported on the current operating system")
}

func (client *SerialClient) ReadByte() (byte, error) {
	panic("not supported on the current operating system")
}
