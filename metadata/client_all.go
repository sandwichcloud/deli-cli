// +build !darwin

package metadata

import (
	"time"

	"github.com/tarm/serial"
)

type SerialClient struct {
	serialClient *serial.Port
}

func (client *SerialClient) Connect(serialPort string) error {
	config := &serial.Config{
		Name:        serialPort,
		Baud:        115200,
		ReadTimeout: 1 * time.Second,
	}

	serialClient, err := serial.OpenPort(config)
	if err != nil {
		return nil
	}
	client.serialClient = serialClient
	return nil
}

func (client *SerialClient) WriteLine(line string) error {
	_, err := client.serialClient.Write([]byte(line + "\n"))
	return err
}

func (client *SerialClient) ReadByte() (byte, error) {
	buf := make([]byte, 1)

	n, _ := client.serialClient.Read(buf)

	if n == 0 {
		return byte(0), NothingToRead
	}

	return buf[0], nil
}
