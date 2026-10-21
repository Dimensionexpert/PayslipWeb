package whatsapp

import "testing"

func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "already international",
			input: "+91 98765 43210",
			want:  "919876543210",
		},
		{
			name:  "with spaces and dashes",
			input: "98765-43210",
			want:  "919876543210",
		},
		{
			name:  "with leading zero",
			input: "09876543210",
			want:  "919876543210",
		},
		{
			name:  "ten digit local",
			input: "9876543210",
			want:  "919876543210",
		},
		{
			name:  "international without plus",
			input: "919876543210",
			want:  "919876543210",
		},
		{
			name:  "international with 00 prefix",
			input: "00919876543210",
			want:  "919876543210",
		},
		{
			name:  "with parentheses",
			input: "(98765) 43210",
			want:  "919876543210",
		},
		{
			name:  "with dots",
			input: "98765.43210",
			want:  "919876543210",
		},
		{
			name:  "with surrounding whitespace",
			input: "  9876543210  ",
			want:  "919876543210",
		},
		{
			name:    "empty",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "only whitespace",
			input:   "   ",
			want:    "",
			wantErr: true,
		},
		{
			name:    "too short",
			input:   "12345",
			want:    "",
			wantErr: true,
		},
		{
			name:    "too long",
			input:   "9876543210123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "contains letters",
			input:   "98765abc210",
			want:    "",
			wantErr: true,
		},
		{
			name:    "contains misplaced plus",
			input:   "91+9876543210",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid starting digit",
			input:   "1234567890",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizePhoneNumber(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizePhoneNumber(%q) error = %v, wantErr %v",
					tt.input, err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("NormalizePhoneNumber(%q) = %q, want %q",
					tt.input, got, tt.want)
			}
		})
	}
}
