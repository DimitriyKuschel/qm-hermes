package interfaces

type CompressorInterface interface {
	Compress(val []byte) ([]byte, error)
	Decompress(val []byte) ([]byte, error)
}
