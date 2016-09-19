package reflux

import (
	"testing"
)

func TestParseInfluxUrlWithUserPass(t *testing.T) {
	s, h, u, p, d, err := parseInfluxUrl("http://u.u:brf@influx.host:1234/database")
	if err != nil {
		t.Error(err)
	}
	if s != "http" {
		t.Error("scheme should be 'http' got", s)
	}
	if h != "influx.host:1234" {
		t.Error("host should be 'influx.host:1234' got", h)
	}
	if u != "u.u" {
		t.Error("user should be 'u.u' got", u)
	}
	if p != "brf" {
		t.Error("password should be 'brf' got", p)
	}
	if d != "database" {
		t.Error("database name should be 'database' got", d)
	}
}
func TestParseInfluxUrlWithUserOnly(t *testing.T) {
	s, h, u, p, d, err := parseInfluxUrl("http://u.u@influx.host:1234/database")
	if err != nil {
		t.Error(err)
	}
	if s != "http" {
		t.Error("scheme should be 'http' got", s)
	}
	if h != "influx.host:1234" {
		t.Error("host should be 'influx.host:1234' got", h)
	}
	if u != "u.u" {
		t.Error("user should be 'u.u' got", u)
	}
	if p != "" {
		t.Error("password should be '' got", p)
	}
	if d != "database" {
		t.Error("database name should be 'database' got", d)
	}
}
func TestParseInfluxUrl(t *testing.T) {
	s, h, u, p, d, err := parseInfluxUrl("http://influx.host:1234/database")
	if err != nil {
		t.Error(err)
	}
	if s != "http" {
		t.Error("scheme should be 'http' got", s)
	}
	if h != "influx.host:1234" {
		t.Error("host should be 'influx.host:1234' got", h)
	}
	if u != "" {
		t.Error("user should be '' got", u)
	}
	if p != "" {
		t.Error("password should be '' got", p)
	}
	if d != "database" {
		t.Error("database name should be 'database' got", d)
	}
}
