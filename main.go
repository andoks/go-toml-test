package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/davecgh/go-spew/spew"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type config struct {
	General struct {
		NodeName string `toml:"node_name"`
		Location string
	}
	Sensors map[string]struct {
		Name    string
		Channel map[string]struct {
			Address    int64
			SampleFreq int64 `toml:"sample_freq"`
		}
	}
}

var _ fmt.Stringer = config{}

func (c config) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "nodename: %v\n", c.General.NodeName)
	fmt.Fprintf(&buf, "location: %v\n", c.General.Location)

	fmt.Fprint(&buf, "Sensors:")
	for name, sensor := range c.Sensors {
		fmt.Fprintf(&buf, "   sensor \"%v\" (key: %v)\n", sensor.Name, name)
		for chanName, channel := range sensor.Channel {
			fmt.Fprintf(&buf, "       channel \"%v\"\n", chanName)
			fmt.Fprintf(&buf, "           address     : %v\n", channel.Address)
			fmt.Fprintf(&buf, "           sample_freq : %v\n", channel.SampleFreq)
		}
	}

	return buf.String()
}

func main() {
	tomlFile, err := os.Open("test.toml")
	check(err)
	defer tomlFile.Close()

	tomlContent, err := ioutil.ReadAll(tomlFile)
	check(err)

	fmt.Printf("toml-file content:\n----\n%v\n----\n", string(tomlContent))

	// FOR DEBUG PURPOSES TO SEE WHAT TOML LIBRARY ACTUALLY PARSES
	//    // Restart read operation on file
	//    _, err = tomlFile.Seek(0, 0)
	//    check(err)
	//
	//    var vIf interface{}
	//    _, err = toml.DecodeReader(tomlFile, &vIf)
	//    check(err)
	//
	//    fmt.Printf("toml-struct type:\n----\n%T\n----\ncontent:\n----\n%+v\n----\n", vIf, spew.Sdump(vIf))

	// Restart read operation on file
	_, err = tomlFile.Seek(0, 0)
	check(err)

	var value config
	_, err = toml.DecodeReader(tomlFile, &value)
	check(err)

	fmt.Printf("toml-struct type:\n----\n%T\n----\ncontent:\n----\n%+v\n----\n", value, value)
}
