package ports

type Hasher interface {
	Hash(plainPassword string) (string, error)
	Verify(hashedPassword, plainPassword string) bool
}
