package fcm

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"github.com/sukso96100/covid19-push/database"
	"google.golang.org/api/option"
)

type FCMObject struct {
	App       *firebase.App
	MsgClient *messaging.Client
	Ctx       context.Context
}

var fcmApp *FCMObject

func InitFCMApp(credential string) {
	*fcmApp = FCMObject{}
	fcmApp.Init(credential)
}

func GetFCMApp() *FCMObject {
	return fcmApp
}

func (fcm *FCMObject) Init(credential string) {
	fcm.Ctx = context.Background()
	opt := option.WithCredentialsFile(credential)
	app, err := firebase.NewApp(fcm.Ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	fcm.App = app

	client, err := fcm.App.Messaging(fcm.Ctx)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	fcm.MsgClient = client
}

func (fcm *FCMObject) PushStatData(statData database.StatData,
	confirmedInc int, curedInc int, deathInc int) {
	// See documentation on defining a message payload.
	tmpl := "확진:%d(증감:%d), 완치:%d(증감:%d), 사망:%d(증감:%d)"
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "코로나19 확산 현황",
			Body: fmt.Sprintf(tmpl,
				statData.Confirmed, confirmedInc,
				statData.Cured, curedInc,
				statData.Death, deathInc),
		},
		Topic: "stat",
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := fcm.MsgClient.Send(fcm.Ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully sent stat message:", response)

}
