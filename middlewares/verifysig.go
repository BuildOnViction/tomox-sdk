package middlewares

import (
	"net/http"
)

func VerifySignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: Logic to verify the signature of user here
		//if r.Header["Signature"] != nil && r.Header["Hash"] != nil && r.Header["Pubkey"] != nil {
		//
		//	signature := []byte(r.Header["Signature"][0])
		//	addressHash := []byte(r.Header["Hash"][0])
		//	publicKey := crypto.PublicKey(r.Header["Pubkey"][0])
		//	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		//
		//	utils.Logger.Debug(signature)
		//	utils.Logger.Debug(addressHash)
		//	utils.Logger.Debug(publicKeyECDSA)
		//
		//	if !ok {
		//		w.WriteHeader(http.StatusInternalServerError)
		//		return
		//	}
		//
		//	publicKeyBytes := ecrypto.FromECDSAPub(publicKeyECDSA)
		//
		//	sigPublicKey, err := ecrypto.Ecrecover(addressHash, signature)
		//
		//	if err != nil {
		//		w.WriteHeader(http.StatusInternalServerError)
		//		return
		//	}
		//
		//	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
		//	fmt.Println(matches) // true
		//} else {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}
