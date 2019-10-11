package swlib

import (
	"testing"

	"github.com/mdlayher/genetlink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mdlayher/netlink"
)

type connectorMock struct {
	mock.Mock
}

func (m connectorMock) Close() error {
	return nil
}

func (m connectorMock) Execute(msg genetlink.Message, family uint16, flags netlink.HeaderFlags) ([]genetlink.Message, error) {
	args := m.Called(msg, family, flags)
	return args.Get(0).([]genetlink.Message), args.Error(1)
}

func TestConn_ListSwitches(t *testing.T) {
	assert2 := assert.New(t)
	require2 := require.New(t)

	m := &connectorMock{}

	req := genetlink.Message{
		Header: genetlink.Header{
			Command: uint8(CmdGetSwitch),
			Version: 1,
		},
	}
	m.On("Execute", req, uint16(1), netlink.Request|netlink.Dump).Return(
		[]genetlink.Message{
			{
				Header: genetlink.Header{
					Command: 2,
					Version: 1,
				},
				Data: []byte{
					0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x03, 0x00, 0x73, 0x77, 0x69, 0x74,
					0x63, 0x68, 0x30, 0x00, 0x0b, 0x00, 0x04, 0x00, 0x6d, 0x74, 0x37, 0x35, 0x33, 0x30, 0x00, 0x00,
					0x0b, 0x00, 0x05, 0x00, 0x6d, 0x74, 0x37, 0x35, 0x33, 0x30, 0x00, 0x00, 0x08, 0x00, 0x06, 0x00,
					0xff, 0x0f, 0x00, 0x00, 0x08, 0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x00, 0x08, 0x00, 0x09, 0x00,
					0x06, 0x00, 0x00, 0x00, 0x24, 0x00, 0x08, 0x00, 0x04, 0x00, 0x07, 0x00, 0x04, 0x00, 0x07, 0x00,
					0x04, 0x00, 0x07, 0x00, 0x04, 0x00, 0x07, 0x00, 0x04, 0x00, 0x07, 0x00, 0x04, 0x00, 0x07, 0x00,
					0x04, 0x00, 0x07, 0x00, 0x04, 0x00, 0x07, 0x00,
				},
			},
		},
		nil,
	)

	c := Conn{
		conn: m,
		family: genetlink.Family{
			ID:      1,
			Version: 1,
		},
	}

	switches, err := c.ListSwitches()
	require2.NoError(err)

	require2.Len(switches, 1)
	assert2.Equal(Device{
		ID:         1,
		DeviceName: "switch0",
		Alias:      "mt7530",
		Name:       "mt7530",
		VLANs:      4095,
		Ports:      8,
		CPUPort:    6,
		PortMap:    []*PortMap{{1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}, {8, "", 0}},
	}, switches[0])

	m.AssertExpectations(t)
}
