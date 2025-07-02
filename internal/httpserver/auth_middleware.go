package httpserver

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) AuthMiddleware(ctx *gin.Context) {
	signature := ctx.GetHeader("X-Signature-Ed25519")
	if signature == "" {
		ctx.AbortWithStatusJSON(401, nil)
		return
	}

	timestamp := ctx.GetHeader("X-Signature-Timestamp")
	if timestamp == "" {
		ctx.AbortWithStatusJSON(401, errorJson("Missing signature timestamp"))
		return
	}

	// Read the body but make sure it can be consumed again
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		_ = ctx.AbortWithError(500, errors.Wrap(err, "Failed to read body"))
		return
	}

	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Verify signature
	pubKey, err := hex.DecodeString(s.config.Discord.PublicKey)
	if err != nil {
		_ = ctx.AbortWithError(500, errors.Wrap(err, "Failed to decode public key"))
		return
	}

	signatureDecoded, err := hex.DecodeString(signature)
	if err != nil {
		ctx.AbortWithStatusJSON(400, errorJson("Failed to decode signature"))
		return
	}

	payload := append([]byte(timestamp), body...)
	if !ed25519.Verify(pubKey, payload, signatureDecoded) {
		ctx.AbortWithStatusJSON(401, errorJson("Invalid signature"))
		return
	}

	ctx.Next()
}

func errorJson(message string) gin.H {
	return gin.H{
		"error": message,
	}
}
