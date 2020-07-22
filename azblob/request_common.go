package azblob

// ClientProvidedKeyOptions contains headers which may be be specified from service version 2019-02-02 
// or higher to encrypts the data on the service-side with the given key. Use of customer-provided keys 
// must be done over HTTPS. As the encryption key itself is provided in the request, a secure connection 
// must be established to transfer the key. 
type ClientProvidedKeyOptions struct {
	EncryptionKey       *string                  // A Base64-encoded AES-256 encryption key value.
	EncryptionKeySha256 *string                  // The Base64-encoded SHA256 of the encryption key.
	EncryptionAlgorithm EncryptionAlgorithmType // Specifies the algorithm to use when encrypting data using the given key. Must be AES256.
	EncryptionScope     *string
}

// InitClientProvidedKeyOptions function.
// By default the value of encryption algorithm params is "AES256" for service version 2019-02-02 or higher.
func InitClientProvidedKeyOptions(ek *string, eksha256 *string, es *string) (cpk *ClientProvidedKeyOptions) {
	cpk = new(ClientProvidedKeyOptions)
	cpk.EncryptionKey, cpk.EncryptionKeySha256, cpk.EncryptionAlgorithm, cpk.EncryptionScope = ek, eksha256, EncryptionAlgorithmAES256, es
	return cpk
}
