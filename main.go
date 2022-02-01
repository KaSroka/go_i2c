package main

import (
	"log"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	disp_fill_on := [1024 + 1]byte{0x40}
	disp_fill_off := [1024 + 1]byte{0x40}
	for i := 1; i < len(disp_fill_on); i++ {
		disp_fill_on[i] = 0xff
	}

	for _, ref := range i2creg.All() {
		log.Println(ref.Name)
		log.Println(ref.Aliases)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("I2C1")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	log.Print("Sending expander setup... ")
	exp := i2c.Dev{Bus: b, Addr: 0x21}
	_, err = exp.Write([]byte{0x03, 0xfc})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("OK")

	disp := i2c.Dev{Bus: b, Addr: 0x3d}

	init1 := []byte{
		0x00,
		0xAB,       // Set Vdd Mode
		0x01,       // ELW2106AA VCI = 3.0V
		0xAD, 0x9E, // Set IREF selection
		0x15, 0x00, 0x7F, // Set column address
		0x75, 0x00, 0x3F, // Set row address
		0xA0, 0x43, // Set Re-map
		0xA1, 0x00, // Set display start line
		0xA2, 0x00, // Set display offset
		0xA4,       // Set display mode
		0xA8, 0x3F, // Set multiplex ratio
		0xB1, 0x11, // Set Phase1,2 length
		0xB3, 0xF0, // Set display clock divide ratio
		0xB9,       // Grey scale table
		0xBC, 0x04, // Set pre-charge voltage
		0xBE, 0x05, // Set VCOMH deselect level, 0.82 * Vcc
	}

	init2 := []byte{
		0x00,
		0x81, 0x7f, // Set contrast
	}

	init3 := []byte{
		0x00,
		0xaf, // Display on
	}

	cursor := []byte{
		0x00,
		0x15, 0x00, 0x7F, // Set column address
		0x75, 0x00, 0x3F, // Set row address
	}

	log.Print("Sending init1... ")
	_, err = disp.Write(init1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")

	log.Print("Sending init2... ")
	_, err = disp.Write(init2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")

	log.Print("Sending init3... ")
	_, err = disp.Write(init3)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")

	log.Print("Sending cursor... ")
	_, err = disp.Write(cursor)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")

	for {
		log.Print("Sending disp_fill_on...")
		for it := 0; it < 8; it++ {
			_, err = disp.Write(disp_fill_on[:])
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println("OK")

		log.Print("Sending disp_fill_off...")
		for it := 0; it < 8; it++ {
			_, err = disp.Write(disp_fill_off[:])
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println("OK")
	}
}
