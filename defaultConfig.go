package iguana

import (
	"iguana/files"
	"iguana/obfuscator"
	"iguana/utils"
	"log"
)

//GetDefaultConfig returns the default configuration to use to hadle the files
func GetDefaultConfig() Config {
	return Config{
		SWPath:          "sw",
		FuncCheckSum:    utils.GenerateChecksum256,
		FuncIDURLNormal: idURLGenetarorNormal,
		FuncIDURLObf:    idURLGenetarorObf,
		FuncObf:         ObfuscateStatic,
	}
}

func idURLGenetarorNormal(static files.Static) (string, string) {
	id := "/static/" + static.Content.Checksum + static.Extension
	return id, id
}
func idURLGenetarorObf(static files.Static) (string, string) {
	id := "/static/" + static.ContentObf.Checksum + static.Extension
	return id, id
}

func skippingDefault(filePath, pathToSkip string) bool {
	return false
}

//ObfuscateStatic helper to obfuscate statics
func ObfuscateStatic(static files.Static) (string, error) {
	if static.Extension == ".js" || static.Extension == ".mjs" {
		code, err := obfuscator.JS(static.Content.Me, "http://localhost:3000", static.Content.Checksum[1:4]+"_")
		if err != nil {
			log.Printf("iguana -> ObfuscateStatic:1 -> file: %s err: %s", static.Path, err)
			return "", err
		}
		return code, err
	}
	return "", nil
}
