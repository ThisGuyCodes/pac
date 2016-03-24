package pac

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strconv"
	"text/template"
)

// Pac is a http handler that can serve pac files
type Pac struct {
	address  string
	port     int
	template *template.Template
	tld      string
	host     string
}

// Address returns the address to use in a pac file for a given request
func (p Pac) Address() string {
	if p.address != "" {
		return p.address
	}
	host, _, _ := net.SplitHostPort(p.host)
	return host
}

// Port returns the port to use in a pac file for a given request
func (p Pac) Port() int {
	if p.port != 0 {
		return p.port
	}

	_, port, _ := net.SplitHostPort(p.host)
	portNum, _ := strconv.Atoi(port)
	if portNum == 0 {
		return 80
	}
	return portNum
}

// Tld exports the internal tld field
func (p Pac) Tld() string {
	return p.tld
}

// Render renders the underlying template to a string
func (p Pac) Render(host string) string {
	buff := &bytes.Buffer{}
	p.RenderTo(host, buff)

	return buff.String()
}

// RenderTo renders the underlying template to a writer
func (p Pac) RenderTo(host string, wr io.Writer) {
	p.host = host
	p.template.Execute(wr, p)
}

// ServeHTTP is the handler that renders a pac
func (p *Pac) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	p.RenderTo(r.Host, rw)
}

// New creates a new pac http handler
func New(t *template.Template, addr string, port int, tld string) http.Handler {
	return &Pac{
		template: t,
		address:  addr,
		port:     port,
		tld:      tld,
	}
}

// Me is the most basic pac file
var Me = template.Must(template.New("me").Parse(`
function FindProxyForURL (url, host) {
    if (dnsDomainIs(host, '{{.Tld}}')) {
        return 'PROXY {{.Address}}:{{.Port}}'
    }
    return 'DIRECT'
}
`))
