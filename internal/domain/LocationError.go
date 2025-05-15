package domain

type InvalidZipCodeError struct {
	ZipCode string
}

type AddressNotFoundError struct {
	ZipCode string
}

func (e *InvalidZipCodeError) Error() string {
	return "invalid zipcode"
}

func (e *AddressNotFoundError) Error() string {
	return "can not find zipcode"
}
