package queues

import (
	"github.com/klauspost/compress/zstd"
	"queue-manager/internal/queues/interfaces"
)

type ZstdCompression struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func (z *ZstdCompression) Compress(val []byte) ([]byte, error) {
	return z.encoder.EncodeAll(val, make([]byte, 0, len(val))), nil
}

func (z *ZstdCompression) Decompress(val []byte) ([]byte, error) {
	return z.decoder.DecodeAll(val, nil)
}

func NewZstdCompressor() interfaces.CompressorInterface {
	encoder, _ := zstd.NewWriter(nil)
	decoder, _ := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
	return &ZstdCompression{encoder: encoder, decoder: decoder}
}
