package files

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

type staticInfo struct {
	Me              []byte
	Checksum        string
	URL             string
	ID              string
	CheckSumChanged bool
}
