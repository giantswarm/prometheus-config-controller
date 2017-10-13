package certificate

// certificateFile represents a PKI certificate on disk.
type certificateFile struct {
	// path is the location of the certificate on the filesystem.
	path string
	// data is the contents of the certificate.
	data string
}
