/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package smart

import (
	"encoding/hex"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
)

func CheckSign(pub, data, sign string) (bool, error) {
	pk, err := hex.DecodeString(pub)
	if err != nil {
		return false, err
	}
	s, err := hex.DecodeString(sign)
	if err != nil {
		return false, err
	}
	pk = crypto.CutPub(pk)
	return crypto.Verify(pk, []byte(data), s)
}
