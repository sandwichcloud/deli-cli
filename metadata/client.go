package metadata

import (
	"context"
	"encoding/base64"
	"strings"
	"time"
)

type MetaDataClient struct {
	serialClient serialClientInterface
}

type serialClientInterface interface {
	Connect(serialPort string) error
	WriteLine(line string) error
	ReadByte() (byte, error)
}

func (client *MetaDataClient) Connect(serialPort string) error {
	client.serialClient = &SerialClient{}
	return client.serialClient.Connect(serialPort)
}

func (client *MetaDataClient) readLine(ctx context.Context) (string, error) {
	var line string

	for {
		if ctx.Err() != nil {
			break
		}
		b, err := client.serialClient.ReadByte()
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

func (client *MetaDataClient) writePacket(packetCode, data string) error {
	packet := "!!" + packetCode + "#" + base64.StdEncoding.EncodeToString([]byte(data))
	return client.serialClient.WriteLine(packet)
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
