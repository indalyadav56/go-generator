package dto

type EncryptedRequest struct {
	Data     string `json:"data"`
	CrcValue string `json:"crc_value"`
}

type DeviceInfo struct {
	UserID    string
	OS        string
	OSVersion string
	DeviceIP  string
	LatLong   string
}
