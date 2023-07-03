package hash

var _ Hash = (*hash)(nil)

type Conf struct {
	Secret string `toml:"secret"`
	Length int    `toml:"length"`
}
type Hash interface {
	i()

	// HashidsEncode 加密
	HashidsEncode(params []int) (string, error)

	// HashidsDecode 解密
	HashidsDecode(hash string) ([]int, error)
}

type hash struct {
	secret string
	length int
}

func New(secret string, length int) Hash {
	return &hash{
		secret: secret,
		length: length,
	}
}

func (h *hash) i() {}
