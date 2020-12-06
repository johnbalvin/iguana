package obfuscator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"iguana/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func JS(code []byte, url string, prefix string) (string, error) {
	fmt.Print("<")
	data := Obfuscator{Code: string(code)}
	dataOptions := ObfuscatorOptions{
		Compact:                     true,
		DeadCodeInjection:           true,
		DeadCodeInjectionThreshold:  0.5,
		IdentifierNamesGenerator:    "hexadecimal",
		IdentifiersPrefix:           "a_" + prefix, //this is usefull when multiple files are obfuscated, as each global variable get's it's own prefix also it's necesary to avoid variables starting with number in which case it will give you an error
		StringArray:                 true,
		RotateStringArray:           true,
		RotateStringArrayEnabled:    true,
		SelfDefending:               true,
		StringArrayThreshold:        0.8,
		StringArrayThresholdEnabled: true,
		StringArrayEncoding:         "false",
		Target:                      "browser",
		TransformObjectKeys:         true,
		//defaults  can't remove it because json decode slices to nil
		DomainLock:      []string{},
		ReservedNames:   []string{},
		ReservedStrings: []string{},
	}
	data.Options = dataOptions
	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Println("iguana -> obfuscateJS:1 -> err:", err)
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("iguana -> obfuscateJS:2 -> err:", err)
		return "", ErrObfuscatorServer
	}
	switch resp.StatusCode {
	case 200:
		var resultado ObfuscatorAnswer
		if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
			log.Println("iguana -> obfuscateJS:3 -> err:", err)
			return "", err
		}
		return resultado.Code, nil
	case 400: //wrong parameters
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("iguana -> obfuscateJS:4 -> err:", err)
			return "", err
		}
		log.Printf("iguana -> obfuscateJS:5 -> err: %s\n", bodyBytes)
		return "", errors.New(string(bodyBytes))
	case 409: //problem with file itself
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("iguana -> obfuscateJS:5 -> err:", err)
			return "", err
		}
		log.Printf("iguana -> obfuscateJS:6 -> err: %s\n", bodyBytes)
		return "", errors.New(string(bodyBytes))
	default:
		log.Println("iguana -> obfuscateJS:7 -> StatusCode: ", resp.StatusCode)
		return "", utils.ErrUnknown
	}
}
