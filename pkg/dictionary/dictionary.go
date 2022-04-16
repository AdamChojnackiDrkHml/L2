package dictionary

const AlphabetSize = 256

type Dictionary struct {
	dict map[string]int
}

func Dictionary_CreateDictionary() *Dictionary {

	d := &Dictionary{dict: make(map[string]int)}

	for i := 0; i < AlphabetSize; i++ {
		d.dict[string([]byte{byte(i)})] = i
	}

	return d
}

func (d *Dictionary) Dictionary_IsContained(bytes []byte) bool {
	_, ok := d.dict[string(bytes)]

	return ok
}

func (d *Dictionary) Dictionary_AddKey(bytes []byte) {
	d.dict[string(bytes)] = len(d.dict)
}

func (d *Dictionary) Dictionary_GetVal(bytes []byte) int {
	return d.dict[string(bytes)]
}
