package dictionary

import (
	"errorMsgs"
)

type Dictionary map[string]string

func (d Dictionary) Search(word string) (string, error) {
	value, isExist := d[word]

	if isExist {
		return value, nil
	}

	return "", errorMsgs.DictNotFound
}

func (d Dictionary) Add(word, def string) error {
	value, _ := d.Search(word)

	if value != "" {
		return errorMsgs.DictAlreadyRegistered
	}

	d[word] = def

	return nil
}

func (d Dictionary) Update(word, newDef string) error {
	_, err := d.Search(word)

	switch err {
	case nil:
		d[word] = newDef
		return nil
	case errorMsgs.DictNotFound:
		return err
	default:
		return nil
	}
}

func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)

	switch err {
	case nil:
		delete(d, word)
		return nil
	case errorMsgs.DictNotFound:
		return err
	default:
		return nil
	}
}
