package iguana

import (
	"errors"
	"log"
	"regexp"
)

var errDepth = errors.New("Deph over the maximum allow: 30, this would mean you have a circular depencie, or have more than 30 nested dependencies")
var errUnknown = errors.New("The obfuscator throwed an unknown error")
var errObfuscatorServer = errors.New("The obfuscator is not active")

//['"][\w\.\-_/ ]+?((\.js)|(\.json)|(\.woff)|(\.woff2)|(\.ttf)|(\.mjs)|(\.txt)|(\.css)|(\.md)|(\.gif)|(\.ico)|(\.jpeg)|(\.jpg)|(\.png)|(\.tiff)|(\.tif)|(\.svg))['"]`) //matchear todos las url de js,css,mjs y svg para que sean reemplazados
//regJs: regex to find .js files,regServiceW: regex to find service workers definitions
//regStatic: regex to find static files urls, files ending in .js,.json,.woff,.woff2,.ttf,.json,.mjs,.txt,.css,.md,.gif,.ico,.jpeg,.jpg, .png, .tiff, .tif, .svg
var regStatic, regServiceW, regJs, regURL *regexp.Regexp

//FuncObf represents function that is gonna be called to obfuscate the file
type FuncObf func(Static) (string, error)

//FuncString represents will be used for returning id and the url
type FuncString func(Static) (string, string)

//FuncSkipPath represents a func to skip files, first argument should be the file's full path that contains the path to skip
type FuncSkipPath func(string, string) bool

//HTML contains all info about html like if there is any service worker
type HTML struct {
	Path            string
	Content         []byte
	Checksum        string
	ServiceWorkers  map[string]bool
	DependsFullPath map[string]bool
	DataGenerate    bool
}

//Static contains info about the static file like the checksum, its content, and so on
type Static struct {
	Path            string //it's full path
	Name            string
	ChangeContent   bool
	Extension       string
	MimeType        string
	Obfuscate       bool
	Content         staticInfo
	ContentObf      staticInfo
	DependsFullPath map[string]bool
}

//SW service worker
type SW struct {
	FileCaller string //the file url that calls the service workers
	Static
}

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

type staticInfo struct {
	Me              []byte
	Checksum        string
	URL             string
	ID              string
	CheckSumChanged bool
}

type obfuscatorOptions struct {
	Compact                        bool     `json:"compact"`
	ControlFlowFlattening          bool     `json:"controlFlowFlattening"`
	ControlFlowFlatteningThreshold float32  `json:"controlFlowFlatteningThreshold"`
	DeadCodeInjection              bool     `json:"deadCodeInjection"`
	DeadCodeInjectionThreshold     float32  `json:"deadCodeInjectionThreshold"`
	DebugProtection                bool     `json:"debugProtection"`
	DebugProtectionInterval        bool     `json:"debugProtectionInterval"`
	DisableConsoleOutput           bool     `json:"disableConsoleOutput"`
	DomainLock                     []string `json:"domainLock"`
	IdentifierNamesGenerator       string   `json:"identifierNamesGenerator"`
	IdentifiersPrefix              string   `json:"identifiersPrefix"`
	RenameGlobals                  bool     `json:"renameGlobals"`
	ReservedNames                  []string `json:"reservedNames"`
	ReservedStrings                []string `json:"reservedStrings"`
	RotateStringArray              bool     `json:"rotateStringArray"`
	RotateStringArrayEnabled       bool     `json:"rotateStringArrayEnabled"`
	Seed                           int      `json:"seed"`
	SelfDefending                  bool     `json:"selfDefending"`
	SourceMap                      bool     `json:"sourceMap"`
	SourceMapBaseURL               string   `json:"sourceMapBaseUrl"`
	SourceMapFileName              string   `json:"sourceMapFileName"`
	SourceMapMode                  string   `json:"sourceMapMode"`
	SourceMapSeparate              bool     `json:"sourceMapSeparate"`
	StringArray                    bool     `json:"stringArray"`
	StringArrayEncoding            string   `json:"stringArrayEncoding"`
	StringArrayEncodingEnabled     bool     `json:"stringArrayEncodingEnabled"`
	StringArrayThreshold           float32  `json:"stringArrayThreshold"`
	StringArrayThresholdEnabled    bool     `json:"stringArrayThresholdEnabled"`
	Target                         string   `json:"target"`
	TransformObjectKeys            bool     `json:"transformObjectKeys"`
	UnicodeEscapeSequence          bool     `json:"unicodeEscapeSequence"`
}
type obfuscator struct {
	Code    string            `json:"code"`
	Options obfuscatorOptions `json:"options"`
}
type obfuscatorAnswer struct {
	Code      string `json:"code"`
	SourceMap string `json:"sourceMap"`
}

var mimeTypeAllow = map[string]string{
	"js":    "application/javascript",
	"json":  "application/json",
	"woff":  "font/woff",
	"woff2": "font/woff2",
	"ttf":   "font/ttf",
	"mjs":   "text/javascript",
	"html":  "text/html",
	"txt":   "text/plain",
	"css":   "text/css",
	"csv":   "text/csv",
	"md":    "text/markdown",
	"gif":   "image/gif",
	"ico":   "image/x-icon",
	"jpeg":  "image/jpeg",
	"jpg":   "image/jpeg",
	"png":   "image/png",
	"tiff":  "image/tiff",
	"tif":   "image/tiff",
	"svg":   "image/svg+xml",
}
var mimeTypeReplace = map[string]bool{ //mime types that could contains dependencies
	"js":   true,
	"json": true,
	"mjs":  true,
	"css":  true,
}

func init() {
	var err error
	regStatic, err = regexp.Compile(`['"][\w\.\-_/ ]+?((\.js)|(\.json)|(\.woff)|(\.woff2)|(\.ttf)|(\.json)|(\.mjs)|(\.txt)|(\.css)|(\.md)|(\.gif)|(\.ico)|(\.jpeg)|(\.jpg)|(\.png)|(\.tiff)|(\.tif)|(\.svg))['"]`) //matchear todos las url de js,css,mjs y svg para que sean reemplazados
	if err != nil {
		log.Fatal(" err: ", err)
	}
	regServiceW, err = regexp.Compile(`(serviceWorker\.register\()['"][\w\-\._/]+?(\.js)['"]`)
	if err != nil {
		log.Fatal(" err: ", err)
	}
	regJs, err = regexp.Compile(`['"][\w\.\-_/]+?(\.js)['"]`)
	if err != nil {
		log.Fatal(" err: ", err)
	}
	regURL, err = regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(" err: ", err)
	}
}
