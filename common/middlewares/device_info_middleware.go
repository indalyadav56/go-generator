package middlewares

import (
	"common/dto"
	"common/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeviceInfoMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("user_id").(string)

		os := ctx.GetHeader("X-OS")
		if os == "" {
			response.HandleErrorWithAbort(ctx, http.StatusBadRequest, "Missing OS header", nil)
			return
		}

		osVersion := ctx.GetHeader("X-OS-Version")
		if osVersion == "" {
			response.HandleErrorWithAbort(ctx, http.StatusBadRequest, "Missing OS version header", nil)
			return
		}

		deviceIP := ctx.GetHeader("X-Device-Ip")
		if deviceIP == "" {
			response.HandleErrorWithAbort(ctx, http.StatusBadRequest, "Missing Device IP header", nil)
			return
		}

		latLong := ctx.GetHeader("X-Lat-Long")
		if latLong == "" {
			response.HandleErrorWithAbort(ctx, http.StatusBadRequest, "Missing Latitude/Longitude header", nil)
			return
		}

		req := &dto.DeviceInfo{
			UserID:    userID,
			OS:        os,
			OSVersion: osVersion,
			DeviceIP:  deviceIP,
			LatLong:   latLong,
		}

		ctx.Set("device_info", req)
		ctx.Next()
	}
}
