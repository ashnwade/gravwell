/*************************************************************************
 * Copyright 2023 Gravwell, Inc. All rights reserved.
 * Contact: <legal@gravwell.io>
 *
 * This software may be modified and distributed under the terms of the
 * BSD 2-clause license. See the LICENSE file for details.
 **************************************************************************/

package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gravwell/gravwell/v3/ingest"
	"github.com/gravwell/gravwell/v3/ingest/entry"
)

type hecIgst struct {
	GeneratorConfig
	name  string
	auth  string
	uri   *url.URL
	to    time.Duration
	src   net.IP
	tags  map[entry.EntryTag]string
	wg    *sync.WaitGroup
	errch chan error
	wtr   io.WriteCloser
}

func newHecConn(name string, gc GeneratorConfig, to time.Duration) (hec *hecIgst, err error) {
	var uri *url.URL
	if uri, err = url.Parse(gc.HEC); err != nil {
		return
	}
	hec = &hecIgst{
		GeneratorConfig: gc,
		to:              to,
		uri:             uri,
		name:            name,
		tags:            map[entry.EntryTag]string{0: gc.Tag},
		auth:            fmt.Sprintf(`Splunk %s`, gc.Auth),
		errch:           make(chan error, 1),
	}
	if hec.src, err = hec.test(); err != nil {
		return
	}
	rdr, wtr := io.Pipe()
	go hec.httpRoutine(rdr)
	hec.wtr = wtr
	return
}

func (hec *hecIgst) test() (ip net.IP, err error) {
	var conn *net.TCPConn
	var raddr *net.TCPAddr
	var ipstr string

	if raddr, err = net.ResolveTCPAddr(`tcp`, hec.uri.Host); err != nil {
		return
	} else if conn, err = net.DialTCP(`tcp`, nil, raddr); err != nil {
		return
	} else if ipstr, _, err = net.SplitHostPort(conn.LocalAddr().String()); err != nil {
		return
	}
	hec.src = net.ParseIP(ipstr)
	ip = hec.src
	err = conn.Close()
	return
}

func (hec *hecIgst) httpRoutine(rdr io.Reader) {
	var err error
	var req *http.Request
	var resp *http.Response
	var cli http.Client

	defer close(hec.errch)
	if req, err = http.NewRequest(http.MethodPost, hec.uri.String(), rdr); err != nil {
		hec.errch <- err
		return
	}
	req.Header.Add(`Authorization`, hec.auth)
	req.Header.Set(`User-Agent`, hec.name)

	if hec.modeHECRaw {
		//attach URL parameters
		uri := req.URL
		values := uri.Query()
		values.Add(`sourcetype`, hec.Tag)
		req.URL.RawQuery = values.Encode()
	}

	if resp, err = cli.Do(req); err != nil {
		hec.errch <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var msg string
		lr := &io.LimitedReader{R: resp.Body, N: 512}
		if body, err := ioutil.ReadAll(lr); err == nil || err == io.EOF {
			msg = string(body)
		} else {
			fmt.Println("bad read", err)
		}
		hec.errch <- fmt.Errorf("invalid status %d (%v)\n%s", resp.StatusCode, resp.Status, msg)
	} else {
		hec.errch <- nil
	}
	return
}

func (hec *hecIgst) WaitForHot(time.Duration) (err error) {
	return
}

func (hec *hecIgst) Close() (err error) {
	hec.wtr.Close()
	err = <-hec.errch
	return
}

func (hec *hecIgst) Sync(time.Duration) (err error) {
	return //no...
}

func (hec *hecIgst) SourceIP() (net.IP, error) {
	return hec.src, nil
}

func (hec *hecIgst) LookupTag(tag entry.EntryTag) (string, bool) {
	v, ok := hec.tags[tag]
	return v, ok
}

func (hec *hecIgst) NegotiateTag(v string) (tag entry.EntryTag, err error) {
	if err = ingest.CheckTag(v); err != nil {
		return
	}
	for k, vv := range hec.tags {
		if v == vv {
			tag = k
			return
		}
	}

	tag = entry.EntryTag(len(hec.tags))
	hec.tags[tag] = v
	return
}

func (hec *hecIgst) GetTag(v string) (tag entry.EntryTag, err error) {
	for k, vv := range hec.tags {
		if v == vv {
			tag = k
			return
		}
	}
	err = errors.New("not found")
	return
}

func (hec *hecIgst) Write(ts entry.Timestamp, tag entry.EntryTag, data []byte) error {
	return hec.WriteEntry(&entry.Entry{
		TS:   ts,
		Tag:  tag,
		Data: data,
	})
}

func (hec *hecIgst) WriteBatch(ents []*entry.Entry) error {
	for _, v := range ents {
		if err := hec.WriteEntry(v); err != nil {
			return err
		}
	}
	return nil
}

func (hec *hecIgst) WriteEntry(ent *entry.Entry) (err error) {
	if hec.modeHECRaw {
		err = hec.sendRaw(ent)
	} else {
		err = hec.sendEvent(ent)
	}
	return
}

type hecEnt struct {
	Time  float64
	ST    string
	Event json.RawMessage
}

func (hec *hecIgst) sendRaw(ent *entry.Entry) error {
	if _, err := hec.wtr.Write(ent.Data); err != nil {
		return err
	} else if _, err = hec.wtr.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

type hecent struct {
	Event json.RawMessage `json:"event,omitempty"`
	Time  float64         `json:"time,omitempty"`
	ST    string          `json:"sourcetype,omitempty"`
}

var osc bool

func setData(data []byte) json.RawMessage {
	//check if this is a JSON object
	if len(data) >= 2 && data[0] == '{' && data[len(data)-1] == '}' {
		if osc = !osc; osc {
			//return it as is
			return json.RawMessage(data)
		}
		//else fall through and encode as a string
	}
	if v, err := json.Marshal(string(data)); err == nil {
		return json.RawMessage(v)
	}
	return nil
}

func (hec *hecIgst) sendEvent(ent *entry.Entry) (err error) {
	if ent != nil {
		v := hecent{
			Time:  timeFloat(ent.TS),
			Event: json.RawMessage(ent.Data),
			ST:    hec.Tag,
		}
		err = json.NewEncoder(hec.wtr).Encode(v)
	}
	return
}

const (
	TS_SIZE int = 12

	secondsPerMinute       = 60
	secondsPerHour         = 60 * 60
	secondsPerDay          = 24 * secondsPerHour
	secondsPerWeek         = 7 * secondsPerDay
	daysPer400Years        = 365*400 + 97
	daysPer100Years        = 365*100 + 24
	daysPer4Years          = 365*4 + 1
	unixToInternal   int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
)

func timeFloat(ts entry.Timestamp) (r float64) {
	//assign unix time seconds, the HEC ingester can't handle timestamps before the unix EPOC, so if its before that just set it to zero
	if ts.Sec < unixToInternal {
		return //just return zero
	}
	r = float64(ts.Sec - unixToInternal)
	//no add in the nano portion as ms
	r += float64(ts.Nsec) / float64(1000000000)
	return
}
