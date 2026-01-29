package configparse

// Validate performs syntactic validation on raw config
func Validate(raw *RawConfig) error {
	return validateRawSyntax(raw)
}
