package pac

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
)

func TestPacAddress(t *testing.T) {
	t.Parallel()
	this := &Pac{
		address: "donald.trump",
		host:    "something.else:2016",
	}
	if got := this.Address(); got != "donald.trump" {
		t.Fatalf("Bad address, got: %q expected %q\n", got, "donald.trump")
	}
}

func TestPacNoAddress(t *testing.T) {
	t.Parallel()
	this := &Pac{
		host: "also.trump:2016",
	}

	if got := this.Address(); got != "also.trump" {
		t.Fatalf("Bad address, got: %q expected %q\n", got, "donald.trump")
	}
}

func TestPacPort(t *testing.T) {
	t.Parallel()
	this := &Pac{
		port: 2016,
		host: "theywill.payforit:2020",
	}
	if got := this.Port(); got != 2016 {
		t.Fatalf("Bad port, got %d, expected %d\n", got, 2016)
	}
}

func TestPacNoPort(t *testing.T) {
	t.Parallel()
	this := &Pac{
		host: "theywill.payforit:2020",
	}
	if got := this.Port(); got != 2020 {
		t.Fatalf("Bad port, got %d, expected %d\n", got, 2020)
	}
}

func TestPacRenderTo(t *testing.T) {
	t.Parallel()

	this := &Pac{
		address:  "trump.towers",
		port:     2016,
		template: template.Must(template.New("some").Parse("{{.Port}}")),
	}

	buff := &bytes.Buffer{}

	this.RenderTo("", buff)

	if buff.String() != "2016" {
		t.Fatalf("Template rendered badly, got %q expected %q\n", buff.String(), "2016")
	}
}

func TestPacRender(t *testing.T) {
	t.Parallel()

	this := &Pac{
		template: template.Must(template.New("some").Parse("{{.Address}}")),
	}

	str := this.Render("trump.towers:2020")

	if str != "trump.towers" {
		t.Fatalf("Template rendered badly, got %q expected %q\n", str, "trump.towers")
	}
}

func TestPacTld(t *testing.T) {
	t.Parallel()

	this := &Pac{
		tld: "drumpf",
	}

	if str := this.Tld(); str != "drumpf" {
		t.Fatalf("Got bad tld, got %q expected %q\n", str, "drumpf")
	}
}

func TestPacServeHTTP(t *testing.T) {
	t.Parallel()

	this := &Pac{
		tld:      "drumpf",
		template: template.Must(template.New("some").Parse("{{.Tld}}")),
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "bling.drumpf/pac", nil)

	this.ServeHTTP(rec, req)

	if str := rec.Body.String(); str != "drumpf" {
		t.Fatalf("Template rendered badly, got %q expected %q\n", str, "trump.towers")
	}
}
