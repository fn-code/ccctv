package lamcctv

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fn-code/lamcctv/pkg/bot"
)

const (
	CCTVLimitReq = 10
	TimeLimit    = 60 * time.Second
)

type CCTV struct {
	ID   string
	Name string
	URL  string
}

// ReadCCTV contain fetch status
type ReadCCTV struct {
	ID       string
	Name     string
	URL      string
	Status   bool
	Err      error
	NumError int64
}

type CCTVSetup struct {
	Data []*CCTV
	Bot  *bot.BotSetup
}

func New(data []*CCTV, bot *bot.BotSetup) *CCTVSetup {
	return &CCTVSetup{data, bot}
}

func (ct *CCTVSetup) ProsesCCTV() {
	ch := make(chan *ReadCCTV)
	done := make(chan *ReadCCTV)
	errCCTV := make(map[string]*ReadCCTV)

	for _, v := range ct.Data {
		go checkCCTVConnection(v, ch)
	}
	go readCCTVConnectionStatus(ch, done)

	for v := range done {
		if v.Status {
			ctv := &CCTV{
				ID:   v.ID,
				Name: v.Name,
				URL:  v.URL,
			}
			tic := time.NewTimer(TimeLimit)
			go func(ctv *CCTV) {
				<-tic.C
				tic.Stop()
				checkCCTVConnection(ctv, ch)
			}(ctv)

		} else {
			errCtv, ok := errCCTV[v.ID]
			ctv := &CCTV{
				ID:   v.ID,
				Name: v.Name,
				URL:  v.URL,
			}
			switch ok {
			case false:
				v.NumError++
				errCCTV[v.ID] = v
				tic := time.NewTimer(TimeLimit)
				go func(ctv *CCTV) {
					<-tic.C
					tic.Stop()
					checkCCTVConnection(ctv, ch)
				}(ctv)
			case true:
				errCtv.NumError++
				if errCtv.NumError == CCTVLimitReq {
					// send notifications
					fmt.Println("------------------------------Sending Notifiction to ", errCtv.Name, " ---------------------------------")
					errCtv.NumError = 0
					msg := fmt.Sprintf("CCTV %s Bermasalah Tolong Di Periksa Ya", ctv.Name)
					err := ct.SendNotification(msg)
					if err != nil {
						log.Println("Failed Sending Notification ", err)
					}
					tic := time.NewTimer(TimeLimit)
					go func(ctv *CCTV) {
						<-tic.C
						checkCCTVConnection(ctv, ch)
						tic.Stop()
					}(ctv)
					break
				}
				tic := time.NewTimer(TimeLimit)
				go func(ctv *CCTV) {
					<-tic.C
					tic.Stop()
					checkCCTVConnection(ctv, ch)
				}(ctv)

			}
		}
	}
}

func (ct *CCTVSetup) SendNotification(text string) error {
	_, err := ct.Bot.SendMessage(text)
	if err != nil {
		return err
	}
	return nil
}

func checkCCTVConnection(cctv *CCTV, ch chan<- *ReadCCTV) {
	client := http.Client{}
	resp, err := client.Get(cctv.URL)
	if err != nil {
		ch <- &ReadCCTV{cctv.ID, cctv.Name, cctv.URL, false, err, 0}
		return
	}
	ch <- &ReadCCTV{cctv.ID, cctv.Name, cctv.URL, true, nil, 0}
	resp.Body.Close()
}

func readCCTVConnectionStatus(ch <-chan *ReadCCTV, done chan<- *ReadCCTV) {
	for {
		select {
		case res := <-ch:
			if !res.Status {
				fmt.Printf("Failed retrive data from %v\n", res.Name)
				done <- res
				break
			}
			fmt.Printf("Succes retrive data from : %v\n", res.Name)
			done <- res
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
}
