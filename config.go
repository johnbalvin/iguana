package iguana

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
)

func idURLGenetarorNormal(static Static) (string, string) {
	id := "/static/" + static.Content.Checksum + "." + static.Extension
	return id, id
}
func idURLGenetarorObf(static Static) (string, string) {
	id := "/static/" + static.ContentObf.Checksum + "." + static.Extension
	return id, id
}

func skippingDefault(filePath, pathToSkip string) bool {
	return false
}

//ObfuscateStatic helper to obfuscate statics
func ObfuscateStatic(static Static) (string, error) {
	if static.Extension == "js" || static.Extension == "mjs" {
		code, err := obfuscateJS(static.Content.Me, "http://localhost:3000", static.Content.Checksum[1:4]+"_")
		if err != nil {
			log.Printf("iguana -> ObfuscateStatic:1 -> file: %s err: %s", static.Path, err)
			return "", err
		}
		return code, err
	}
	return "", nil
}

//GenerateChecksum256 generates a 256 checksum
func GenerateChecksum256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

//GetDefaultConfig returns the default configuration to use to hadle the files
func GetDefaultConfig() Config {
	return Config{
		SWPath:          "sw",
		FuncCheckSum:    GenerateChecksum256,
		FuncIDURLNormal: idURLGenetarorNormal,
		FuncIDURLObf:    idURLGenetarorObf,
		FuncObf:         ObfuscateStatic,
	}
}
