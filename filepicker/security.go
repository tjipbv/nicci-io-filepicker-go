package filepicker

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// PolicyOpts structure defines the properties of the policy.
type PolicyOpts struct {
	// Expiry is the expiration date of the policy after which it will no longer
	// be valid. This field is required.
	Expiry time.Time `json:"-"`

	// Handle is an unique file handle that you would like to access. It can be
	// obtained from Blob object by calling Handle() method.
	Handle string `json:"handle,omitempty"`

	// Call defines the list of function calls which this policy will be allowed
	// to make.
	Call []Method `json:"call,omitempty"`

	// MaxSize sets the maximum object size limit. This property only applies to
	// the store command.
	MaxSize uint64 `json:"maxsize,omitempty"`

	// MaxSize sets the minimum object size limit. This property only applies to
	// the store command.
	MinSize uint64 `json:"minsize,omitempty"`

	// Path field is valid only for policies that store files. It is a perl-like
	// regular expression that must match the path that the files will be stored
	// under. Defaults to allowing any path ('.*').
	Path string `json:"path,omitempty"`

	// Container field is valid only for policies that store files. It is a
	// perl-like regular expression that must match the container that the files
	// will be stored under. Defaults to allowing any container ('.*').
	Container string `json:"container,omitempty"`
}

// Method defines the calls that created policy will be able to make.
type Method string

// TODO : (ppknap)
const (
	MetPick     = Method("pick")     // Pick methods.
	MetRead     = Method("read")     // Download methods.
	MetStat     = Method("stat")     // Stat method.
	MetWrite    = Method("write")    // Write method.
	MetWriteurl = Method("writeUrl") // WriteURL method.
	MetStore    = Method("store")    // Store methods.
	MetConvert  = Method("convert")  // Convert method.
	MetRemove   = Method("remove")   // Remove method.
)

// MarshalJSON implements json.Marshaler interface. It transforms Expiry field
// representation to UNIX time value. By default, marshaling time.Time structure
// produces a quoted string in RFC 3339 format.
func (po *PolicyOpts) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PolicyOpts
		ExpiryUNIX int64 `json:"expiry"`
	}{*po, po.Expiry.Unix()})
}

// Policy stores the information about what the user can or cannot do.
type Policy string

// MakePolicy creates a new Policy object from provided policy options.
func MakePolicy(po *PolicyOpts) (policy Policy, err error) {
	if po == nil || po.Expiry.IsZero() {
		return policy, fmt.Errorf("filepicker: invalid expiration date")
	}
	byted, err := json.Marshal(po)
	if err != nil {
		return
	}
	return Policy(base64.URLEncoding.EncodeToString(byted)), nil
}

// Security type stores the piece of information that is required to access
// secured URLs.
type Security struct {
	Policy    Policy `json:"policy,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// MakeSecurity creates a new Security object from the given secret and policy
// instances.
//
// You should not store your secret in your code. Instead, call this function
// once and then use obtained strings to initialize Security objects directly.
func MakeSecurity(secret string, policy Policy) Security {
	hasher := hmac.New(sha256.New, []byte(secret))
	hasher.Write([]byte(policy))
	return Security{
		Policy:    policy,
		Signature: hex.EncodeToString(hasher.Sum(nil)),
	}
}
