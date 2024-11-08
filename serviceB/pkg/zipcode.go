package pkg

import "regexp"

func IsValidZipCode(cep string) bool {
	cepRegex := regexp.MustCompile(`^\d{5}-?\d{3}$`)

	if cepRegex.MatchString(cep) {
		return true
	}

	return false
}
