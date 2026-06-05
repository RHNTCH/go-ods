package ods

// Close closes the underlying content stream and zip reader.
func (r *Reader) Close() error {
	var err error

	if r.content != nil {
		err = r.content.Close()
	}

	if r.zipReader != nil {
		if zipErr := r.zipReader.Close(); zipErr != nil {
			return zipErr
		}
	}

	return err
}
