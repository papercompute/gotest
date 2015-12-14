package xlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//    "log"
	"errors"
	"sort"
	"strings"
)

type xTime struct {
	time.Time
}

func (t *xTime) UnmarshalJSON(buf []byte) error {
	tm, err := time.Parse("2006-01-02", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tm
	return nil
}

type Period struct {
	Start xTime `json:"start"`
	End   xTime `json:"end"`
}

func (p *Period) String() string {
	return fmt.Sprintf("{\"start\": \"%s\", \"end\": \"%s\"}",
		p.Start.Time.Format("2006-01-02"), p.End.Time.Format("2006-01-02"))
}

type Data struct {
	RentingPeriod Period   `json:"renting_period"`
	Reservations  []Period `json:"reservations"`
}

type Periods []Period

func (a Periods) Len() int           { return len(a) }
func (a Periods) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Periods) Less(i, j int) bool { return compDate(a[i].Start, a[j].Start) < 0 }

func ReadJSON(r *http.Request, v interface{}) error {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, v); err != nil {
		return err
	}
	return nil
}

func compDate(a, b xTime) int {
	ay, am, ad := a.Time.Date()
	by, bm, bd := b.Time.Date()
	if ay > by {
		return 1
	}
	if ay < by {
		return -1
	}

	if am > bm {
		return 1
	}
	if am < bm {
		return -1
	}

	if ad > bd {
		return 1
	}
	if ad < bd {
		return -1
	}

	return 0
}

func LongestFreeRange(d *Data) (string, error) {
	if len(d.Reservations) == 0 {
		return d.RentingPeriod.String(), nil
	}

	start := d.RentingPeriod.Start
	end := d.RentingPeriod.End

	if compDate(start, end) > 0 {
		return "", errors.New(fmt.Sprintf("invalid renting period %v - %v", start, end))
	}

	var ps []Period

	sort.Sort(Periods(d.Reservations))

	for i, t := range d.Reservations {
		// check start <= end
		if compDate(t.Start, t.End) > 0 {
			return "", errors.New(fmt.Sprintf("invalid interval %v", t))
		}
		// check reservations overlap
		if i > 0 {
			u := d.Reservations[i-1]
			if compDate(t.Start, u.End) <= 0 {
				return "", errors.New(fmt.Sprintf("reservations overlap %v,%v", u, t))
			}
		}

		if compDate(t.Start, start) >= 0 {
			if compDate(t.Start, end) > 0 {
				continue
			} else {
				if compDate(t.End, end) > 0 {
					ps = append(ps, Period{Start: t.Start, End: end})
				} else {
					ps = append(ps, t)
				}
			}
		} else {
			if compDate(t.End, start) < 0 {
				continue
			} else {
				if compDate(t.End, end) > 0 {
					ps = append(ps, Period{Start: start, End: end})
				} else {
					ps = append(ps, Period{Start: start, End: t.End})
				}

			}

		}
	}

	var ps2 []Period

	if compDate(start, ps[0].Start) < 0 {
		ps2 = append(ps2, Period{start, xTime{ps[0].Start.AddDate(0, 0, -1)}})
	}
	l := len(ps) - 1
	if compDate(ps[l].End, end) < 0 {
		ps2 = append(ps2, Period{xTime{ps[l].End.AddDate(0, 0, 1)}, end})
	}

	for i := 0; i < l; i++ {
		t := Period{xTime{ps[i].End.AddDate(0, 0, 1)},
			xTime{ps[i+1].Start.AddDate(0, 0, -1)}}

		if (int)(t.End.Time.Sub(t.Start.Time)/(time.Hour*24)) >= 0 {
			ps2 = append(ps2, t)
		}
	}

	maxDura := ps2[0].End.Time.Sub(ps2[0].Start.Time)
	maxIdx := 0
	for i := 1; i < len(ps2); i++ {
		d := ps2[i].End.Time.Sub(ps2[i].Start.Time)
		//log.Printf("[%d] %v : %v",i,d,maxDura)
		if d > maxDura {
			maxDura = d
			maxIdx = i
		}
	}

	/*
		for i,t:=range ps2{
			s:=fmt.Sprintf("%s - %s : %d",
				t.Start.Time.Format("2006-01-02"),
				t.End.Time.Format("2006-01-02"),
				(int)(t.End.Time.Sub(t.Start.Time)/(time.Hour*24)))
				if i==maxIdx{
					s=s+"*"
				}
			log.Printf(s)
		}
	*/
	return ps2[maxIdx].String(), nil
}
