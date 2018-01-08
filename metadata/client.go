package metadata

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type MetaDataClient struct {
	SerialPort   string
	serialClient *serial.Port
}

func (client *MetaDataClient) Connect() error {
	config := &serial.Config{
		Name:        client.SerialPort,
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

func (client *MetaDataClient) readByte() (byte, error) {
	buf := make([]byte, 1)

	n, _ := client.serialClient.Read(buf)

	if n == 0 {
		return byte(0), NothingToRead
	}

	return buf[0], nil
}

func (client *MetaDataClient) readLine(ctx context.Context) (string, error) {
	var line string

	for {
		if ctx.Err() != nil {
			break
		}
		b, err := client.readByte()
		if err != nil {
			if err == NothingToRead {
				continue
			}
			return "", err
		}
		line += string(b)

		if strings.HasSuffix(line, "\n") {
			break
		}
	}

	return line, nil
}

func (client *MetaDataClient) readPacket(ctx context.Context) (string, string, error) {
	line, err := client.readLine(ctx)
	if err != nil {
		return "", "", err
	}

	if strings.HasPrefix(line, "!!") == false {
		return "", "", InvalidPacketError
	}

	packet := strings.Split(line, "#")
	packetCode := packet[0][2:]
	packetData := packet[1]

	packetDataBytes, err := base64.StdEncoding.DecodeString(packetData)
	if err != nil {
		return "", "", err
	}

	return packetCode, string(packetDataBytes), nil
}

func (client *MetaDataClient) writeLine(line string) error {
	_, err := client.serialClient.Write([]byte(line + "\n"))
	return err
}

func (client *MetaDataClient) writePacket(packetCode, data string) error {
	packet := "!!" + packetCode + "#" + base64.StdEncoding.EncodeToString([]byte(data))
	return client.writeLine(packet)
}

func (client *MetaDataClient) getData(txPacketCode, rxPacketCode string) (string, error) {
	c := make(chan string)
	defer close(c)

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	go func() {
		var packetData string
		var packetCode string

		for packetCode != rxPacketCode && err == nil {
			packetCode, packetData, err = client.readPacket(ctx)
			if err == InvalidPacketError { // Packet had invalid data so keep reading
				err = nil
				continue
			}
		}

		if err != nil {
			cancel()
			return
		}

		c <- packetData
		ctx.Done()
	}()

	err = client.writePacket(txPacketCode, "")
	if err != nil {
		cancel()
	}

	var output string
	select {
	case <-ctx.Done():
		if ctx.Err() != nil {
			if ctx.Err() == context.Canceled {
				return "", err
			}
			return "", TimedOut
		}
		return "", ctx.Err()
	case output = <-c:
		break
	}

	return output, nil
}

func (client *MetaDataClient) GetMetaData() (string, error) {
	return client.getData("REQUEST_METADATA", "RESPONSE_METADATA")
}

func (client *MetaDataClient) GetNetworkData() (string, error) {
	return client.getData("REQUEST_NETWORKDATA", "RESPONSE_NETWORKDATA")
}

func (client *MetaDataClient) GetUserData() (string, error) {
	return client.getData("REQUEST_USERDATA", "RESPONSE_USERDATA")
}

func (client *MetaDataClient) GetSecurityData() (string, error) {
	return client.getData("REQUEST_SECURITYDATA", "RESPONSE_SECURITYDATA")
}
