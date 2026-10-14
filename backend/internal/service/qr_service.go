package service

import (
	"fmt"
	"github.com/skip2/go-qrcode"
)

type QRService struct{}

func NewQRService() *QRService {
	return &QRService{}
}

// GeneratePNG generates a PNG QR code byte slice for the given target URL
func (s *QRService) GeneratePNG(content string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256
	}
	if size > 1024 {
		size = 1024
	}

	png, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate qr code: %w", err)
	}

	return png, nil
}
