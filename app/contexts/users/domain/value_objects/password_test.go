package value_objects

import (
	"errors"
	"testing"

	domain_exceptions "github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeHasher struct {
	hashReturn      string
	verifyReturn    bool
	hashErr         error
	lastHashInput   string
	lastVerifyHash  string
	lastVerifyPlain string
}

func (f *fakeHasher) Hash(plain string) (string, error) {
	f.lastHashInput = plain
	return f.hashReturn, f.hashErr
}

func (f *fakeHasher) Verify(hash, plain string) bool {
	f.lastVerifyHash = hash
	f.lastVerifyPlain = plain
	return f.verifyReturn
}

func TestPasswordVerify_ReturnsTrueAndUsesHasher(t *testing.T) {
	fh := &fakeHasher{
		hashReturn:   "hashed123",
		verifyReturn: true,
	}

	p, err := NewPassword(fh, "Abcdefg1")
	require.NoError(t, err)
	assert.Equal(t, "Abcdefg1", fh.lastHashInput, "Hasher.Hash should receive the plain password")

	ok := p.Verify("Abcdefg1")
	assert.True(t, ok, "expected Verify to return true")
	assert.Equal(t, "hashed123", fh.lastVerifyHash, "expected Verify to use hashed value")
	assert.Equal(t, "Abcdefg1", fh.lastVerifyPlain, "expected Verify to receive plain")
}

func TestPasswordVerify_ReturnsFalseWhenHasherDeclines(t *testing.T) {
	fh := &fakeHasher{
		hashReturn:   "h2",
		verifyReturn: false,
	}

	p, err := NewPassword(fh, "Abcdefg1")
	require.NoError(t, err)

	ok := p.Verify("Wrongpass1")
	assert.False(t, ok, "expected Verify to return false")
	assert.Equal(t, "h2", fh.lastVerifyHash, "expected Verify to use hashed value 'h2'")
	assert.Equal(t, "Wrongpass1", fh.lastVerifyPlain, "expected Verify to receive plain 'Wrongpass1'")
}

func TestNewPassword_SetsHashAndVerifierWorks(t *testing.T) {
	fh := &fakeHasher{
		hashReturn:   "h123",
		verifyReturn: true,
	}

	p, err := NewPassword(fh, "Abcdefg1")
	require.NoError(t, err)
	assert.Equal(t, "h123", p.Hash())

	assert.True(t, p.Verify("Abcdefg1"))
	assert.Equal(t, "h123", fh.lastVerifyHash)
	assert.Equal(t, "Abcdefg1", fh.lastVerifyPlain)
}

func TestNewPassword_WeakPasswords(t *testing.T) {
	cases := []struct {
		name string
		pw   string
	}{
		{name: "too short", pw: "Abc1e"},
		{name: "no uppercase", pw: "abcdefg1"},
		{name: "no lowercase", pw: "ABCDEFG1"},
		{name: "no digit", pw: "Abcdefgh"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fh := &fakeHasher{}
			_, err := NewPassword(fh, tc.pw)
			require.Error(t, err)
			assert.ErrorIs(t, err, domain_exceptions.ErrWeakPassword)
			assert.Empty(t, fh.lastHashInput, "hasher should not be used on weak passwords")
		})
	}
}

func TestNewPassword_HashErrorIsWrapped(t *testing.T) {
	fh := &fakeHasher{
		hashErr: errors.New("boom"),
	}
	_, err := NewPassword(fh, "Abcdefg1")
	require.Error(t, err)
	assert.ErrorContains(t, err, "failed to hash password")
}

func TestNewPasswordFromHash_SetsHash(t *testing.T) {
	p := NewPasswordFromHash("static-hash")
	assert.Equal(t, "static-hash", p.Hash())
}

func TestNewPasswordFromHash_VerifyPanicsWithoutHasher(t *testing.T) {
	p := NewPasswordFromHash("static-hash")
	require.Panics(t, func() {
		_ = p.Verify("anything")
	}, "Verify should panic when hasher is nil")
}
