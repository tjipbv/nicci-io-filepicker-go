package filepicker

import (
	//"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// TODO : (ppknap)
type Method string

const (
	Pick     = Method("pick")
	Read     = Method("read")
	Stat     = Method("stat")
	Write    = Method("write")
	Writeurl = Method("writeUrl")
	Store    = Method("store")
	Convert  = Method("convert")
	Remove   = Method("remove")
)

// TODO : (ppknap)
type PolicyOpts struct {
	Expiry    time.Time `json:"-"`
	Handle    string    `json:"handle,omitempty"`
	Call      []Method  `json:"call,omitempty"`
	MaxSize   uint64    `json:"maxsize,omitempty"`
	MinSize   uint64    `json:"minsize,omitempty"`
	Path      string    `json:"path,omitempty"`
	Container string    `json:"container,omitempty"`
}

// MarshalJSON implements json.Marshaler interface. It transforms Expiry field
// representation to UNIX time value. By default, marshaling time.Time structure
// produces a quoted string in RFC 3339 format.
func (po *PolicyOpts) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PolicyOpts
		ExpiryUNIX int64 `json:"expiry"`
	}{*po, po.Expiry.Unix()})
}

// TODO : (ppknap)
type Policy string

// TODO : (ppknap)
func MakePolicy(po *PolicyOpts) (policy Policy, err error) {
	if po.Expiry.IsZero() {
		return policy, fmt.Errorf("filepicker: invalid expiration date")
	}
	byted, err := json.Marshal(po)
	if err != nil {
		return
	}
	return Policy(base64.URLEncoding.EncodeToString(byted)), nil
}

// TODO : (ppknap)
type Security struct {
	Policy    Policy `json:"policy,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// toValues takes all non-zero values from provided Security instance and puts
// them to a url.Values object.
func (s Security) toValues() url.Values {
	return toValues(s)
}

// TODO : (ppknap)
func MakeSecurity(secret string, policy Policy) (security Security) {
	hasher := hmac.New(sha256.New, []byte(secret))
	hasher.Write([]byte(policy))
	return Security{
		Policy:    policy,
		Signature: hex.EncodeToString(hasher.Sum(nil)),
	}
}
