package pin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	_ "log"
	"strconv"
)

type PinBlock_Iso3 struct {
	PinBlocker
}

func (pin_block *PinBlock_Iso3) Encrypt(pan_12digits string, clear_pin string, key []byte) []byte {

	if len(clear_pin) > 12 {
		panic("pin length > 12")
	}

	buf := bytes.NewBufferString(fmt.Sprintf("3%X%s", len(clear_pin), clear_pin))

	//random pads
	fill_random(buf)

	pin_block_data_a, _ := hex.DecodeString(buf.String())
	pin_block_data_b, _ := hex.DecodeString("0000" + pan_12digits)

	for i, v := range pin_block_data_b {
		pin_block_data_a[i] = pin_block_data_a[i] ^ v
	}
	//log.Printf(" xor'ed pin block =", hex.EncodeToString(pin_block_data_a))

	enc_pin_block := EncryptPinBlock(pin_block_data_a, key)
	return (enc_pin_block)

}

func (pin_block *PinBlock_Iso3) GetPin(pan_12digits string, pin_block_data []byte, key []byte) string {

	clear_pin_block := DecryptPinBlock(pin_block_data, key)

	//pan_12digits := pan[len(pan)-13 : len(pan)-1]
	pin_block_data_b, _ := hex.DecodeString("0000" + pan_12digits)
	//log.Printf(" pin block (b) =", hex.EncodeToString(pin_block_data_b))

	for i, v := range pin_block_data_b {
		clear_pin_block[i] = clear_pin_block[i] ^ v
	}

	pin_block_str := hex.EncodeToString(clear_pin_block)
	//log.Printf(" clear pin block (b) =", pin_block_str)

	n_pin_digits, _ := strconv.ParseInt(pin_block_str[1:2], 16, 16)
	clear_pin := pin_block_str[2:(2 + n_pin_digits)]

	return (clear_pin)

}
