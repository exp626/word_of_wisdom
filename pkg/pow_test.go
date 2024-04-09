package pkg

import "testing"

func Test_PoW(t *testing.T) {
	tests := []struct {
		name    string
		n, k, d int
	}{
		{
			name: "successful pow",
			n:    48,
			k:    2,
			d:    3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewEquihashPoW(tt.n, tt.d, tt.k)

			data, err := client.Challenge()
			if err != nil {
				t.Fatal(err)
			}

			nonce, solution, err := client.PoW(data)
			if err != nil {
				t.Fatal(err)
			}

			if !client.Validate(nonce, solution, data) {
				t.Errorf("Validate() = %v, want %v", false, true)
			}

		})
	}
}
