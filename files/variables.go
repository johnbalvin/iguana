package files

import "regexp"

//['"][\w\.\-_/ ]+?((\.js)|(\.json)|(\.woff)|(\.woff2)|(\.ttf)|(\.mjs)|(\.txt)|(\.css)|(\.md)|(\.gif)|(\.ico)|(\.jpeg)|(\.jpg)|(\.png)|(\.tiff)|(\.tif)|(\.svg))['"]`) //matchear todos las url de js,css,mjs y svg para que sean reemplazados
//regJs: regex to find .js files,regServiceW: regex to find service workers definitions
//regStatic: regex to find static files urls, files ending in .js,.json,.woff,.woff2,.ttf,.json,.mjs,.txt,.css,.md,.gif,.ico,.jpeg,.jpg, .png, .tiff, .tif, .svg
var RegStatic = regexp.MustCompile(`['"][\w\.\-_/ ]+?((\.js)|(\.json)|(\.woff)|(\.woff2)|(\.ttf)|(\.json)|(\.mjs)|(\.txt)|(\.css)|(\.md)|(\.gif)|(\.ico)|(\.jpeg)|(\.jpg)|(\.png)|(\.tiff)|(\.tif)|(\.svg))['"]`) //matchear todos las url de js,css,mjs y svg para que sean reemplazados
var RegServiceW = regexp.MustCompile(`(serviceWorker\.register\()['"][\w\-\._/]+?(\.js)['"]`)
var RegJs = regexp.MustCompile(`['"][\w\.\-_/]+?(\.js)['"]`)
var RegURL = regexp.MustCompile("[^a-zA-Z0-9]+")

var MimeTypeAllow = map[string]string{
	".js":    "application/javascript",
	".json":  "application/json",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".ttf":   "font/ttf",
	".mjs":   "text/javascript",
	".html":  "text/html",
	".txt":   "text/plain",
	".css":   "text/css",
	".csv":   "text/csv",
	".md":    "text/markdown",
	".gif":   "image/gif",
	".ico":   "image/x-icon",
	".jpeg":  "image/jpeg",
	".jpg":   "image/jpeg",
	".png":   "image/png",
	".tiff":  "image/tiff",
	".tif":   "image/tiff",
	".svg":   "image/svg+xml",
}
var MimeTypeReplace = map[string]bool{ //mime types that could contains dependencies
	".js":   true,
	".json": true,
	".mjs":  true,
	".css":  true,
}
