/*
 * Copyright (C) 2020-2022, IrineSistiana
 *
 * This file is part of mosdns.
 *
 * mosdns is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mosdns is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package resp_ip

import (
	"github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/pkg/matcher/netlist"
	"github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/pkg/query_context"
	"github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/sequence"
	"github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/base_ip"
	"github.com/miekg/dns"
	"net"
	"net/netip"
)

const PluginType = "resp_ip"

func init() {
	sequence.MustRegMatchQuickSetup(PluginType, QuickSetup)
}

type Args = base_ip.Args

func QuickSetup(bq sequence.BQ, s string) (sequence.Matcher, error) {
	return base_ip.NewMatcher(bq, base_ip.ParseQuickSetupArgs(s), matchRespAddr)
}

func matchRespAddr(qCtx *query_context.Context, m netlist.Matcher) (bool, error) {
	r := qCtx.R()
	if r == nil {
		return false, nil
	}
	for _, rr := range r.Answer {
		var ip net.IP
		switch rr := rr.(type) {
		case *dns.A:
			ip = rr.A
		case *dns.AAAA:
			ip = rr.AAAA
		default:
			continue
		}
		addr, ok := netip.AddrFromSlice(ip)
		if ok && m.Match(addr) {
			return true, nil
		}
	}
	return false, nil
}
