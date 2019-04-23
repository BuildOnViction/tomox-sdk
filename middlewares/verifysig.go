package middlewares

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tomochain/tomodex/utils"
)

func VerifySignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: Logic to verify the signature of user here
		if r.Header["Signature"] != nil && r.Header["Hash"] != nil && r.Header["Pubkey"] != nil {

			signature := common.Hex2Bytes(r.Header["Signature"][0])
			hash := common.Hex2Bytes(r.Header["Hash"][0])
			publicKeyBytes := common.Hex2Bytes(r.Header["Pubkey"][0])

			utils.Logger.Debug(signature)
			utils.Logger.Debug(hash)
			utils.Logger.Debug(publicKeyBytes)

			sigPublicKey, err := crypto.Ecrecover(hash, signature)
			if err != nil {
				utils.Logger.Error(err)
			}

			matches := bytes.Equal(sigPublicKey, publicKeyBytes)
			fmt.Println(matches) // true

			if !matches {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
