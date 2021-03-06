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

type arrayConfig struct {
	Location            string
	LogLevel            string `toml:"log_level"`
	TagDataWithHostname bool   `toml:"tag_data_with_hostname"`
	Sensors             []struct {
		Name     string
		UUID     string
		Channels []struct {
			Name       string
			Address    int64
			SampleFreq int64 `toml:"sample_freq"`
		}
	}
}

var _ fmt.Stringer = arrayConfig{}

func (c arrayConfig) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "location                 : %v\n", c.Location)
	fmt.Fprintf(&buf, "log-level                : %v\n", c.LogLevel)
	fmt.Fprintf(&buf, "tag data with hostname   : %v\n", c.TagDataWithHostname)

	fmt.Fprint(&buf, "Sensors:\n")
	for _, sensor := range c.Sensors {
		fmt.Fprintf(&buf, "   sensor \"%v\" (UUID: \"%v\")\n", sensor.Name, sensor.UUID)
		for _, channel := range sensor.Channels {
			fmt.Fprintf(&buf, "       channel \"%v\"\n", channel.Name)
			fmt.Fprintf(&buf, "           address     : %v\n", channel.Address)
			fmt.Fprintf(&buf, "           sample_freq : %v\n", channel.SampleFreq)
		}
	}

	return buf.String()
}

type subtableConfig struct {
	Location            string
	LogLevel            string `toml:"log_level"`
	TagDataWithHostname bool   `toml:"tag_data_with_hostname"`
	Sensors             map[string]struct {
		UUID     string
		Channels map[string]struct {
			Address    int64
			SampleFreq int64 `toml:"sample_freq"`
		}
	}
}

var _ fmt.Stringer = subtableConfig{}

func (c subtableConfig) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "location                 : %v\n", c.Location)
	fmt.Fprintf(&buf, "log-level                : %v\n", c.LogLevel)
	fmt.Fprintf(&buf, "tag data with hostname   : %v\n", c.TagDataWithHostname)

	fmt.Fprint(&buf, "Sensors:\n")
	for name, sensor := range c.Sensors {
		fmt.Fprintf(&buf, "   sensor \"%v\" (key: %v)\n", sensor.UUID, name)
		for chanName, channel := range sensor.Channels {
			fmt.Fprintf(&buf, "       channel \"%v\"\n", chanName)
			fmt.Fprintf(&buf, "           address     : %v\n", channel.Address)
			fmt.Fprintf(&buf, "           sample_freq : %v\n", channel.SampleFreq)
		}
	}

	return buf.String()
}

func testTomlParse(filename string, objectToParseInto interface{}) {
	tomlFile, err := os.Open(filename)
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

	_, err = toml.DecodeReader(tomlFile, objectToParseInto)
	check(err)

	objAsStringer := objectToParseInto.(fmt.Stringer)
	fmt.Printf("toml-struct type:\n----\n%T\n----\nstruct content:\n----\n%+v\n----\n", objectToParseInto, objAsStringer)
}

func main() {
	arrConf := arrayConfig{}
	testTomlParse("test-array-of-tables.toml", &arrConf)

	subtabConf := subtableConfig{}
	testTomlParse("test-subtables.toml", &subtabConf)
}
