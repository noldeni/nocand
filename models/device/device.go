package device

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// Info
//
//
type Info struct {
	Signature    [4]byte
	VersionMajor byte
	VersionMinor byte
	ChipId       [12]byte
}

// StatusByte
//
//
type StatusByte byte

const (
	STATUS_RX_PENDING = 0x01
	STATUS_TX_PENDING = 0x02
	STATUS_UNDEF_0    = 0x04
	STATUS_UNDEF_1    = 0x08
	STATUS_ERROR      = 0x10
	STATUS_FAULT      = 0x20
	STATUS_POWERED    = 0x40
	STATUS_CAN_RES    = 0x80
)

var StatusStrings = [...]string{
	"rx-pending",
	"tx-pending",
	"undefined_0",
	"undefined_1",
	"driver error",
	"electric fault",
	"powered",
	"resistor",
}

func (s StatusByte) Strings() []string {
	var i uint
	var r []string

	for i = 0; i < 8; i++ {
		if (s & (1 << i)) != 0 {
			r = append(r, StatusStrings[i])
		} else {
			if i == 6 {
				r = append(r, "unpowered")
			}
		}
	}
	return r
}

func (s StatusByte) String() string {
	var r string
	for _, s := range s.Strings() {
		r += "+" + s
	}
	return r
}

func (s StatusByte) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strings.Join(s.Strings(), ", ") + "\""), nil
}

// PowerStatus
//
//
type PowerStatus struct {
	Status       StatusByte `json:"status"`
	Voltage      float32    `json:"voltage"`
	CurrentSense uint16     `json:"current_sense"`
	RefLevel     float32    `json:"reference_voltage"`
}

func (ps PowerStatus) String() string {
	return fmt.Sprintf("Driver voltage=%.1f, current sense=%d, reference voltage=%.2f, status(%x)=%s.", ps.Voltage, ps.CurrentSense, ps.RefLevel, byte(ps.Status), ps.Status)
}

func (ps *PowerStatus) PackValue() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, ps)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (ps *PowerStatus) UnpackValue(b []byte) error {
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, ps)
	if err != nil {
		return err
	}
	return nil
}
