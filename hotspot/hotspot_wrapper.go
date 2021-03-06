package hotspot

import (
	"fmt"
	"strings"
	"os/exec"
)


type HotspotWrapper struct {
	ConfigCommand string
	ActiveCommand string
	StopCommand string

	SSID string
	Password string
}

func (hotspot *HotspotWrapper) FillInfo(ssid, password string ) {
	hotspot.ConfigCommand = fmt.Sprintf(hotspot.ConfigCommand, ssid, password)
}

func (hotspot *HotspotWrapper) ConfigProperties(ssid, password string ) {

	hotspot.SSID = ssid
	hotspot.Password = password

	if strings.Contains(hotspot.ConfigCommand, "%") {
		hotspot.ConfigCommand = fmt.Sprintf(hotspot.ConfigCommand, ssid, password)
	}
}

func (hotspot *HotspotWrapper) CheckPcSupport () bool {
	checkCommand := "NETSH WLAN show drivers"
	chunks := strings.Split(checkCommand, " ")

	c := exec.Command(chunks[0], chunks[1:]...)

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
		return false
	}

	response, err := c.Output()
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	expectedSpanish := "Red hospedada admitida: sí"
	expectedEnglish := "Hosted network supported: Yes"

	if strings.Contains(string(response), expectedSpanish) ||
		strings.Contains(string(response), expectedEnglish) {
		return true
	}
	return false

}

func (hotspot *HotspotWrapper) SecondConfig() error {
	command1 := "NETSH WLAN set hostednetwork ssid=" + hotspot.SSID
	command2 := "NETSH WLAN set hostednetwork key=" + hotspot.Password

	chunkCommand1 := strings.Split(command1, " ")
	c1 := exec.Command(chunkCommand1[0], chunkCommand1[1:]...)

	if err := c1.Run(); err != nil {
		return err
	}

	chunkCommand2 := strings.Split(command2, " ")
	c2 := exec.Command(chunkCommand2[0], chunkCommand2[1:]...)

	if err := c2.Run(); err != nil {
		return err
	}

	return nil
}

func (hotspot *HotspotWrapper) ActivateHotspot() error {
	chunks := strings.Split(hotspot.ActiveCommand, " ")
	c := exec.Command(chunks[0], chunks[1:]...)
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

func (hotspot *HotspotWrapper) StopHotspot() error{
	chunks := strings.Split(hotspot.StopCommand, " ")
	c := exec.Command(chunks[0], chunks[1:]...)
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}



func NewHotspotWrapper(os string) *HotspotWrapper {
	if os == "windows" {
		return &HotspotWrapper{
			ConfigCommand: "NETSH WLAN set hostednetwork mode=allow ssid=%s key=%s",
			ActiveCommand: "NETSH WLAN start hostednetwork",
			StopCommand: "NETSH WLAN stop hostednetwork",
		}
	}
	return nil
}
