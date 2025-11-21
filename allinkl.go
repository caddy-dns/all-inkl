package allinkl

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	allinkl "github.com/libdns/all-inkl"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *allinkl.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.allinkl",
		New: func() caddy.Module { return &Provider{new(allinkl.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.KasUsername = repl.ReplaceAll(p.Provider.KasUsername, "")
	p.Provider.KasPassword = repl.ReplaceAll(p.Provider.KasPassword, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	allinkl [<username> <password>] {
//	    kas_username <username>
//	    kas_password <password>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.KasUsername = d.Val()
		}
		if d.NextArg() {
			p.Provider.KasPassword = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "kas_username":
				if d.NextArg() {
					p.Provider.KasUsername = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "kas_password":
				if d.NextArg() {
					p.Provider.KasPassword = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.KasUsername == "" {
		return d.Err("missing KasUsername")
	}

	if p.Provider.KasPassword == "" {
		return d.Err("missing KasPassword")
	}

	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
