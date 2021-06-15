package serialize


import (
	"encoding/gob"
	"os"
)

type Serialize struct {
}

func NewSerialize()*Serialize{
	return new(Serialize)
}

func (s *Serialize)SerializeStructToFile(fileName string, t interface{}) error {

	_, err:= os.Stat(fileName)

	var file *os.File

	if err != nil {// not exits
		file, err = os.Create(fileName)
	}

	file, err = os.OpenFile(fileName, os.O_TRUNC | os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(t)

	if err != nil {
		return err
	}
	return nil
}

func (s *Serialize)UnSerializeStructFromFile(fileName string, t interface{})error{
	_, err:= os.Stat(fileName)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}

	defer file.Close()

	dec := gob.NewDecoder(file)
	err = dec.Decode(t)
	if err != nil{
		return err
	}

	return nil
}