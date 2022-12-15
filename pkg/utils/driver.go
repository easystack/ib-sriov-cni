package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	sysBusPciDrivers = "/sys/bus/pci/drivers"
	// UserspaceDrivers is a list of driver names that don't have netlink representation for their devices
	UserspaceDrivers = []string{"vfio-pci", "uio_pci_generic", "igb_uio"}
)

// HasUserspaceDriver checks if a device is attached to Userspace
func HasUserspaceDriver(pciAddr string) (bool, error) {
	driverName, err := GetDriverName(pciAddr)
	if err != nil {
		return false, err
	}
	for _, drv := range UserspaceDrivers {
		if driverName == drv {
			return true, nil
		}
	}
	return false, nil
}

// GetDriverName returns current driver attached to a pci device from its pci address
func GetDriverName(pciAddr string) (string, error) {
	driverLink := filepath.Join(SysBusPci, pciAddr, "driver")
	driverPath, err := filepath.EvalSymlinks(driverLink)
	if err != nil {
		return "", fmt.Errorf("error getting driver info for device %s %v", pciAddr, err)
	}
	driverStat, err := os.Stat(driverPath)
	if err != nil {
		return "", fmt.Errorf("error getting driver stat for device %s %v", pciAddr, err)
	}
	return driverStat.Name(), nil
}

// Unbind unbind driver for one device
func UnbindVf(pciAddr, driver string) error {
	filePath := filepath.Join(sysBusPciDrivers, driver, "unbind")
	err := ioutil.WriteFile(filePath, []byte(pciAddr), os.ModeAppend)
	if err != nil {
		return fmt.Errorf("fail to unbind driver for device %s. %s", pciAddr, err)
	}
	return nil
}

// Bind bind driver for one device
func BindVf(pciAddr, driver string) error {
	filePath := filepath.Join(sysBusPciDrivers, driver, "bind")
	err := ioutil.WriteFile(filePath, []byte(pciAddr), os.ModeAppend)
	if err != nil {
		return fmt.Errorf("fail to bind driver for device %s. %s", pciAddr, err)
	}
	return nil
}
