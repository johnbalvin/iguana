package files

import (
	"bytes"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

func CompressBrotli(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w := brotli.NewWriterLevel(&buf, brotli.BestCompression)
	defer w.Close()

	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CompressZstd(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	enc, err := zstd.NewWriter(
		&buf,
		zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(22)),
		zstd.WithEncoderConcurrency(1),
	)
	if err != nil {
		return nil, err
	}
	defer enc.Close()

	if _, err := enc.Write(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
