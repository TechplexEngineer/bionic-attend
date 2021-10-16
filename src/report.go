package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func (s Server) handleReport() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		type RowData struct {
			FirstName string
			LastName  string
			Total     int64
			Percent   float64
			Meetings  map[string]string
		}

		// map of [userid] to RowData
		Rows := map[string]*RowData{
			//"123": {
			//	FirstName: "Blake",
			//	LastName:  "Bourque",
			//	Total:     2,
			//
			//	Meetings: map[string]string{
			//		"2021-10-1": "x",
			//		"2021-10-2": "",
			//		"2021-10-3": "x",
			//	},
			//},
			//"456": {
			//	FirstName: "Ice",
			//	LastName:  "Pack",
			//	Total:     2,
			//	Meetings: map[string]string{
			//		"2021-10-1": "",
			//		"2021-10-2": "x",
			//		"2021-10-3": "x",
			//	},
			//},
		}

		attendance, err := s.db.GetAttendance(context.Background())
		if err != nil {
			err = fmt.Errorf("error GetAttendance - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}

		meetings, err := s.db.GetMeetings(context.Background())
		if err != nil {
			err = fmt.Errorf("error GetMeetings - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}

		totalMeetings := len(meetings)

		for _, row := range attendance {
			userRow, ok := Rows[row.Userid]
			if !ok { // missing, so we create the entry
				d := RowData{
					FirstName: row.FirstName,
					LastName:  row.LastName,
					Total:     1,
					Percent:   1.0 / float64(totalMeetings),
					Meetings:  make(map[string]string),
				}
				d.Meetings[row.Date] = "x"
				Rows[row.Userid] = &d
				log.Printf("Creating: %#v %v", row.Userid, row.Date)
				continue
			}
			// not missing so update
			userRow.Meetings[row.Date] = "x"
			userRow.Total++
			userRow.Percent = float64(userRow.Total) / float64(totalMeetings)
			log.Printf("Updating: %s %v", row.Userid, row.Date)
		}

		log.Printf("Rows: %#v", Rows)

		d := struct {
			Meetings      []string
			TotalMeetings int
			Rows          map[string]*RowData
		}{
			Meetings:      meetings, //[]string{"2021-10-1", "2021-10-2", "2021-10-3"},
			Rows:          Rows,
			TotalMeetings: totalMeetings,
		}

		// generate data

		//

		//
		//users, err := s.db.ListUsers(context.Background())
		//if err != nil {
		//	err = fmt.Errorf("error ListUsers - %w", err)
		//	s.handleInternalError(err)(w, req)
		//	return
		//}

		//d := struct {
		//	Meetings []string
		//	Users    []data.User
		//}{
		//	Meetings: meetings,
		//	Users:    users,
		//}

		// handle the request
		s.handleTemplate("report.html", d)(w, req)
	}
}
