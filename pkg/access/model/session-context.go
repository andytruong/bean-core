package model

import (
	"fmt"
	"io"
	"strconv"
)

type DeviceType string

const (
	// which source: user/password, sso(microsoft, google, github, facebook, twitter, …)
	// two-step login
	DeviceTypeDesktop    DeviceType = "Desktop"
	DeviceTypeLaptop     DeviceType = "Laptop"
	DeviceTypeSmartPhone DeviceType = "SmartPhone"
	DeviceTypeTablet     DeviceType = "Tablet"
	DeviceTypeTv         DeviceType = "TV"
)

var AllDeviceType = []DeviceType{
	DeviceTypeDesktop,
	DeviceTypeLaptop,
	DeviceTypeSmartPhone,
	DeviceTypeTablet,
	DeviceTypeTv,
}

func (e DeviceType) IsValid() bool {
	switch e {
	case DeviceTypeDesktop, DeviceTypeLaptop, DeviceTypeSmartPhone, DeviceTypeTablet, DeviceTypeTv:
		return true
	}
	return false
}

func (e DeviceType) String() string {
	return string(e)
}

func (e *DeviceType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DeviceType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DeviceType", str)
	}
	return nil
}

func (e DeviceType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
