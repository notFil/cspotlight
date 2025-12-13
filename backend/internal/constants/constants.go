package constants

const (
	UserContextKey = "userContext"
)

var CriticalDirectives = []string{"script-src", "script-src-elem", "script-src-attr", "default-src", "object-src", "worker-src", "trusted-types", "require-trusted-types-for"}
var HighDirectives = []string{"connect-src", "form-action", "frame-src", "frame-ancestors", "child-src", "manifest-src", "base-uri"}
var MediumDirectives = []string{"style-src", "style-src-elem", "style-src-attr", "frame-src", "font-src", "media-src", "img-src", "prefetch-src"}
var LowDirectives = []string{"navigate-to", "sandbox", "upgrade-insecure-requests"}
