package obfuscator

type ObfuscatorOptions struct {
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
type Obfuscator struct {
	Code    string            `json:"code"`
	Options ObfuscatorOptions `json:"options"`
}
type ObfuscatorAnswer struct {
	Code      string `json:"code"`
	SourceMap string `json:"sourceMap"`
}
