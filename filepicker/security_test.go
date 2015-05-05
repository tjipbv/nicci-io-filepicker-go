package filepicker_test

import (
	"testing"
	"time"

	"github.com/filepicker/filepicker-go/filepicker"
)

func TestSecurity(t *testing.T) {
	tests := []struct {
		Key string
		Opt *filepicker.PolicyOpts
		Sec filepicker.Security
	}{
		{
			Key: `Z3IYZSH2UJA7VN3QYFVSVCF7PI`,
			Opt: &filepicker.PolicyOpts{
				Expiry: time.Unix(1508141504, 0),
				Handle: `KW9EJhYtS6y48Whm2S6D`,
			},
			Sec: filepicker.Security{
				Policy:    `eyJoYW5kbGUiOiJLVzlFSmhZdFM2eTQ4V2htMlM2RCIsImV4cGlyeSI6MTUwODE0MTUwNH0=`,
				Signature: `4098f262b9dba23e4766ce127353aaf4f37fde0fd726d164d944e031fd862c18`,
			},
		},
		{
			Key: `S4IXZSH2UJA7VN3QYFVSVCF7PI`,
			Opt: &filepicker.PolicyOpts{
				Expiry: time.Unix(1508154321, 0),
				Call: []filepicker.Method{
					filepicker.MetStore,
					filepicker.MetWrite,
				},
			},
			Sec: filepicker.Security{
				Policy:    `eyJjYWxsIjpbInN0b3JlIiwid3JpdGUiXSwiZXhwaXJ5IjoxNTA4MTU0MzIxfQ==`,
				Signature: `d458b7c957080835815d5394d84b85522e9964dc9853a615a7c3651cf96768e3`,
			},
		},
	}

	for _, test := range tests {
		policy, err := filepicker.MakePolicy(test.Opt)
		if err != nil {
			t.Errorf(`want err == nil; got %v`, err)
		}
		if policy != test.Sec.Policy {
			t.Errorf(`want policy == test.Sec.Policy; got %v != %v`, policy, test.Sec.Policy)
		}
		security := filepicker.MakeSecurity(test.Key, policy)
		if security.Signature != test.Sec.Signature {
			t.Errorf(`want security.Signeture == test.Sec.Signature; got %v != %v`,
				security.Signature, test.Sec.Signature)
		}
	}
}

func TestSecurityErrorNoTime(t *testing.T) {
	_, err := filepicker.MakePolicy(&filepicker.PolicyOpts{
		Handle: `KW9EJhYtS6y48Whm2S6D`,
	})
	if err == nil {
		t.Errorf(`want err != nil; got nil`)
	}
}
