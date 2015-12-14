package main

import (
	"./src"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var d xlib.Data
		if err := xlib.ReadJSON(r, &d); err != nil {
			http.Error(w, fmt.Sprintf("request error %v", err), http.StatusBadRequest)
			return
		}

		//log.Printf("[%s - %s]",
		//	d.RentingPeriod.Start.Time.Format("2006-01-02"),
		//	d.RentingPeriod.End.Time.Format("2006-01-02"))

		for _, t := range d.Reservations {
			log.Printf("%s - %s",
				t.Start.Time.Format("2006-01-02"),
				t.End.Time.Format("2006-01-02"))
		}

		p, err := xlib.LongestFreeRange(&d)
		if err != nil {
			http.Error(w, fmt.Sprintf("request error %v", err), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(p))
		return
	}

	http.Error(w, "request error", http.StatusBadRequest)
}

func main() {
	pp := flag.String("p", "8080", "port")
	flag.Parse()
	http.HandleFunc("/longest-free-range", handler)
	log.Printf("ListenAndServe :" + *pp)
	log.Fatal(http.ListenAndServe(":"+*pp, nil))
}
