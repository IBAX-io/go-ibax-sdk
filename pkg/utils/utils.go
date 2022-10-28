package utils

import (
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/pkg/common/crypto"
	"github.com/IBAX-io/go-ibax-sdk/pkg/consts"
	"github.com/IBAX-io/go-ibax-sdk/pkg/converter"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
)

// CheckSign checks the signature
func CheckSign(publicKeys [][]byte, forSign []byte, signs []byte, nodeKeyOrLogin bool) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{"type": consts.PanicRecoveredError, "error": r}).Error("recovered panic in check sign")
		}
	}()

	var signsSlice [][]byte
	if len(forSign) == 0 {
		log.WithFields(log.Fields{"type": consts.EmptyObject}).Error("for sign is empty")
		return false, ErrInfoFmt("len(forSign) == 0")
	}
	if len(publicKeys) == 0 {
		log.WithFields(log.Fields{"type": consts.EmptyObject}).Error("public keys is empty")
		return false, ErrInfoFmt("len(publicKeys) == 0")
	}
	if len(signs) == 0 {
		log.WithFields(log.Fields{"type": consts.EmptyObject}).Error("signs is empty")
		return false, ErrInfoFmt("len(signs) == 0")
	}

	// node always has only one signature
	if nodeKeyOrLogin {
		signsSlice = append(signsSlice, signs)
	} else {
		length, err := converter.DecodeLength(&signs)
		if err != nil {
			log.WithFields(log.Fields{"type": consts.UnmarshallingError, "error": err}).Error("decoding signs length")
			return false, err
		}
		if length > 0 {
			signsSlice = append(signsSlice, converter.BytesShift(&signs, length))
		}

		if len(publicKeys) != len(signsSlice) {
			log.WithFields(log.Fields{"public_keys_length": len(publicKeys), "signs_length": len(signsSlice), "type": consts.SizeDoesNotMatch}).Error("public keys and signs slices lengths does not match")
			return false, fmt.Errorf("sign error publicKeys length %d != signsSlice length %d", len(publicKeys), len(signsSlice))
		}
	}

	return crypto.Verify(publicKeys[0], forSign, signsSlice[0])
}

// ErrInfoFmt fomats the error message
func ErrInfoFmt(err string, a ...any) error {
	return fmt.Errorf("%s (%s)", fmt.Sprintf(err, a...), Caller(1))
}

// Caller returns the name of the latest function
func Caller(steps int) string {
	name := "?"
	if pc, _, num, ok := runtime.Caller(steps + 1); ok {
		name = fmt.Sprintf("%s :  %d", filepath.Base(runtime.FuncForPC(pc).Name()), num)
	}
	return name
}
