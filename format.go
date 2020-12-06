package iguana

import "iguana/files"

//FuncObf represents function that is gonna be called to obfuscate the file
type FuncObf func(files.Static) (string, error)

//FuncString represents will be used for returning id and the url
type FuncString func(files.Static) (string, string)

//FuncSkipPath represents a func to skip files, first argument should be the file's full path that contains the path to skip
type FuncSkipPath func(string, string) bool

//Config it will have the config to use
type Config struct {
	SWPath             string              //the main path of your sw files
	FuncCheckSum       func([]byte) string //funtion to get checksum
	FuncIDURLNormal    FuncString          //function to get id and url from file if content it's not obfuscated
	FuncIDURLObf       FuncString          //function to get id and url from file if content it's obfuscated
	FuncObf            FuncObf             //function to get the file's obfuscated content
	FuncReplaceRelPath FuncSkipPath        //function to skip files, it returns true, file will be skip for replacement
	SkipLogging        bool
}
