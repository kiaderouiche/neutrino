package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
)

type JSON map[string]interface{}

func (j JSON) ForEach(f func(key string, value interface{})) {
	for k := range j {
		f(k, j[k])
	}
}

func (j JSON) String() string {
	var b bytes.Buffer

	j.ForEach(func(k string, v interface{}) {
		b.WriteString(k + ":" + fmt.Sprintf("%v", v) + "\r\n")
	})

	return b.String()
}

func (j JSON) JsonString() ([]byte, error) {
	return json.Marshal(j)
}

func (j JSON) FromMap(m map[string]interface{}) JSON {
	for k := range m {
		j[k] = m[k]
	}

	return j
}

func (j JSON) FromResponse(res *http.Response) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, j)
	if err != nil {
		return err
	}

	return nil
}