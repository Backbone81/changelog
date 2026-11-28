package changelog

import "errors"

type Severity int

const (
	// To keep backwards compatibility (also visually within the yaml files) we need the default (0) to not have an actual value so it will be omitted when writing the yaml output.
	SeverityUnknown Severity = iota
	SeverityMajor
	SeverityMinor
	SeverityPatch
)

func ParseSeverity(severity string) (Severity, error) {
	switch severity {
	case "major":
		return SeverityMajor, nil
	case "minor":
		return SeverityMinor, nil
	case "patch":
		return SeverityPatch, nil
	default:
		return SeverityUnknown, errors.New("severity must be one of major, minor or patch")
	}
}
