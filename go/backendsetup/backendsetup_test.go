package backendsetup

import (
	"testing"
)

func TestIsIPv6(t *testing.T) {
	tests := []struct {
		addr string
		want bool
	}{
		{"10.0.0.1", false},
		{"192.168.1.100", false},
		{"0.0.0.0", false},
		{"2001:db8::1", true},
		{"fc00::1", true},
		{"::1", true},
		{"not-an-ip", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			got := IsIPv6(tt.addr)
			if got != tt.want {
				t.Errorf("IsIPv6(%q) = %v, want %v", tt.addr, got, tt.want)
			}
		})
	}
}

func TestDeduplicateVIPs(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "no duplicates",
			input: []string{"10.0.0.1", "10.0.0.2", "2001:db8::1"},
			want:  []string{"10.0.0.1", "10.0.0.2", "2001:db8::1"},
		},
		{
			name:  "with duplicates",
			input: []string{"10.0.0.1", "10.0.0.2", "10.0.0.1", "10.0.0.2"},
			want:  []string{"10.0.0.1", "10.0.0.2"},
		},
		{
			name:  "empty list",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "single element",
			input: []string{"10.0.0.1"},
			want:  []string{"10.0.0.1"},
		},
		{
			name:  "all duplicates",
			input: []string{"10.0.0.1", "10.0.0.1", "10.0.0.1"},
			want:  []string{"10.0.0.1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeduplicateVIPs(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("DeduplicateVIPs() returned %d elements, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("DeduplicateVIPs()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestExtractVIPAddresses(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    []string
		wantErr bool
	}{
		{
			name: "basic vips",
			yaml: `
vips:
  - address: "192.168.1.100"
    port: 80
    proto: "tcp"
  - address: "192.168.1.101"
    port: 443
    proto: "tcp"
`,
			want: []string{"192.168.1.100", "192.168.1.101"},
		},
		{
			name: "duplicate addresses",
			yaml: `
vips:
  - address: "192.168.1.100"
    port: 80
    proto: "tcp"
  - address: "192.168.1.100"
    port: 443
    proto: "tcp"
`,
			want: []string{"192.168.1.100"},
		},
		{
			name: "mixed v4 and v6",
			yaml: `
vips:
  - address: "192.168.1.100"
    port: 80
    proto: "tcp"
  - address: "2001:db8::1"
    port: 443
    proto: "tcp"
`,
			want: []string{"192.168.1.100", "2001:db8::1"},
		},
		{
			name: "no vips section",
			yaml: `
server:
  port: 8080
`,
			want: []string{},
		},
		{
			name:    "invalid yaml",
			yaml:    `{{{invalid`,
			wantErr: true,
		},
		{
			name: "full config with other sections",
			yaml: `
server:
  port: 8080
lb:
  interfaces:
    main: "eth0"
target_groups:
  web:
    - address: "10.0.0.1"
      weight: 100
vips:
  - address: "192.168.1.100"
    port: 80
    proto: "tcp"
    target_group: web
`,
			want: []string{"192.168.1.100"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractVIPAddresses([]byte(tt.yaml))
			if (err != nil) != tt.wantErr {
				t.Fatalf("extractVIPAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("extractVIPAddresses() returned %d elements, want %d: %v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("extractVIPAddresses()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuildIPCommand(t *testing.T) {
	tests := []struct {
		name   string
		action string
		vip    string
		want   []string
	}{
		{
			name:   "add ipv4",
			action: "add",
			vip:    "10.0.0.1",
			want:   []string{"addr", "add", "10.0.0.1/32", "dev", "lo"},
		},
		{
			name:   "add ipv6",
			action: "add",
			vip:    "2001:db8::1",
			want:   []string{"addr", "add", "2001:db8::1/128", "dev", "lo"},
		},
		{
			name:   "del ipv4",
			action: "del",
			vip:    "192.168.1.100",
			want:   []string{"addr", "del", "192.168.1.100/32", "dev", "lo"},
		},
		{
			name:   "del ipv6",
			action: "del",
			vip:    "fc00::1",
			want:   []string{"addr", "del", "fc00::1/128", "dev", "lo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildIPCommand(tt.action, tt.vip)
			if len(got) != len(tt.want) {
				t.Fatalf("BuildIPCommand() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("BuildIPCommand()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuildTunnelCommand(t *testing.T) {
	tests := []struct {
		iface string
		want  []string
	}{
		{"tunl0", []string{"link", "set", "tunl0", "up"}},
		{"ip6tnl0", []string{"link", "set", "ip6tnl0", "up"}},
	}
	for _, tt := range tests {
		t.Run(tt.iface, func(t *testing.T) {
			got := BuildTunnelCommand(tt.iface)
			if len(got) != len(tt.want) {
				t.Fatalf("BuildTunnelCommand() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("BuildTunnelCommand()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestDetectTunnelNeeds(t *testing.T) {
	tests := []struct {
		name    string
		vips    []string
		wantV4  bool
		wantV6  bool
	}{
		{
			name:   "only ipv4",
			vips:   []string{"10.0.0.1", "192.168.1.1"},
			wantV4: true, wantV6: false,
		},
		{
			name:   "only ipv6",
			vips:   []string{"2001:db8::1", "fc00::1"},
			wantV4: false, wantV6: true,
		},
		{
			name:   "mixed",
			vips:   []string{"10.0.0.1", "2001:db8::1"},
			wantV4: true, wantV6: true,
		},
		{
			name:   "empty",
			vips:   []string{},
			wantV4: false, wantV6: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV4, gotV6 := detectTunnelNeeds(tt.vips)
			if gotV4 != tt.wantV4 || gotV6 != tt.wantV6 {
				t.Errorf("detectTunnelNeeds() = (%v, %v), want (%v, %v)",
					gotV4, gotV6, tt.wantV4, tt.wantV6)
			}
		})
	}
}
