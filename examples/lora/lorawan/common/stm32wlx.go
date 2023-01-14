//go:build stm32wlx

package common

import (
	"machine"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/sx126x"
)

const (
	FREQ                      = 868100000
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

var (
	loraRadio *sx126x.Device
)

var spi = machine.SPI3

func newRadioControl() sx126x.RadioController {
	return sx126x.NewRadioControl()
}

// do sx126x setup here
func SetupLora() (lora.Radio, error) {
	loraRadio = sx126x.New(spi)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)

	// Create radio controller for target
	loraRadio.SetRadioController(newRadioControl())

	if state := loraRadio.DetectDevice(); !state {
		return nil, errRadioNotFound
	}

	loraConf := lora.Config{
		Freq:           FREQ,
		Bw:             lora.Bandwidth_500_0,
		Sf:             lora.SpreadingFactor9,
		Cr:             lora.CodingRate4_7,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       12,
		Ldr:            lora.LowDataRateOptimizeOff,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       lora.SyncPrivate,
		LoraTxPowerDBm: 20,
	}

	loraRadio.LoraConfig(loraConf)

	return loraRadio, nil
}

func FirmwareVersion() string {
	return "sx126x"
}

func Lorarx() ([]byte, error) {
	return loraRadio.Rx(LORA_DEFAULT_RXTIMEOUT_MS)
}