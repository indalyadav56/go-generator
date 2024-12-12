package utils

import (
	"common/dto"
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromContext(ctx *gin.Context) (string, error) {
	userID, ok := ctx.MustGet("user_id").(string)
	if !ok {
		return "", errors.New("user_id is required")
	}
	return userID, nil
}

func GetRequestIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value("request_id").(string)
	if !ok {
		return ""
	}
	return userID
}

func GetDecryptedDataFromContext(ctx *gin.Context) (string, error) {
	decryptedData, ok := ctx.MustGet("decrypted").(string)
	if !ok {
		return "", errors.New("request decryption payload missing")
	}
	return decryptedData, nil
}

func GetDeviceInfoFromContext(ctx *gin.Context) (*dto.DeviceInfo, error) {
	deviceInfo, ok := ctx.MustGet("device_info").(*dto.DeviceInfo)
	if !ok {
		return nil, errors.New("device info is required")
	}

	return deviceInfo, nil
}
