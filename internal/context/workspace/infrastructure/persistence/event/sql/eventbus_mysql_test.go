package sql

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/matchers"
	"github.com/spf13/viper"

	"github.com/Bio-OS/bioos/internal/apiserver/options"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/db"
	"github.com/Bio-OS/bioos/pkg/log"
)

func TestMysql(t *testing.T) {
	uri := os.Getenv("MYSQL_URI")
	if len(uri) == 0 {
		t.Logf("set env MYSQL_URI to enable TestMongoDB. e.g. root:12345678@tcp(127.0.0.1:3306)/test, you can setup a mysql in veay sample way `docker run --name some-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=12345678 -e MYSQL_DATABASE=test -d mysql`")
		return
	}
	g := gomega.NewWithT(t)
	ctx := context.Background()
	log.RegisterLogger(log.NewOptions())

	conf, err := mysql.ParseDSN(uri)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	host, port, err := net.SplitHostPort(conf.Addr)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	viper.AutomaticEnv()
	_ = os.Setenv("MYSQL_HOST", host)
	_ = os.Setenv("MYSQL_PORT", port)
	_ = os.Setenv("MYSQL_USERNAME", conf.User)
	_ = os.Setenv("MYSQL_PASSWORD", conf.Passwd)
	_ = os.Setenv("MYSQL_DB", conf.DBName)
	opts := &options.Options{
		DBOption: &db.Options{
			MySQL: &db.MySQLOptions{
				Username: "MYSQL_USERNAME",
				Password: "MYSQL_PASSWORD",
				Host:     "MYSQL_HOST",
				Port:     "MYSQL_PORT",
				Database: "MYSQL_DB",
			},
		},
	}

	var eventBus eventbus.EventBus
	var eventRepo eventbus.EventRepository
	orm, err := opts.DBOption.MySQL.GetGORMInstance(ctx)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	eventRepo, err = NewEventRepository(ctx, orm, time.Minute*5, time.Minute*60)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	eOpts := []eventbus.Option{
		eventbus.WithMaxRetries(1),
		eventbus.WithSyncPeriod(time.Second * 5),
		eventbus.WithBatchSize(10),
	}
	eventBus, err = eventbus.NewEventBus(eventRepo, eOpts...)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	go func() {
		err := eventBus.Start(ctx, 1)
		g.Expect(err).ToNot(gomega.HaveOccurred())
	}()

	defer orm.Delete(&Event{})
	orm.Delete(&Event{})

	eventType := "payloadChan"
	payload := []byte("id1")
	payload2 := []byte("id1")

	handler := fakeEventbusOkHandler{
		payloadChan:     make(chan []byte, 1),
		expectedPayload: payload,
	}

	eventBus.Subscribe(eventType, handler)
	failedHandler := fakeEventbusFailedHandler{
		typeA:    make(chan []byte, 1),
		expected: payload,
	}
	eventBus.Subscribe(eventType, failedHandler)

	err = eventBus.Publish(ctx, fakeEvent{
		EventTyp: eventType,
		ID:       string(payload),
	})
	err = eventBus.Publish(ctx, fakeEvent{
		EventTyp: eventType,
		ID:       string(payload2),
	})

	//select {}
	g.Expect(err).ToNot(gomega.HaveOccurred())

	g.Expect(handler.IsOk()).To(gomega.BeTrue())
	g.Expect(failedHandler.IsOk()).To(gomega.BeFalse())
	time.Sleep(time.Second * 1)
	events, err := eventRepo.ListAndLockUnfinishedEvents(ctx, 100, []string{})
	g.Expect(err).ToNot(gomega.HaveOccurred())
	for _, event := range events {
		g.Expect(event.Status).To(&matchers.EqualMatcher{
			Expected: eventbus.EventStatusFailed,
		})
	}
}

type fakeEvent struct {
	EventTyp string
	ID       string
}

func (e fakeEvent) EventType() string {
	return e.EventTyp
}
func (e fakeEvent) Payload() []byte {
	return []byte(e.ID)
}
func (e fakeEvent) Delay() time.Duration {
	return 0
}

type fakeEventbusOkHandler struct {
	payloadChan     chan []byte
	expectedPayload []byte
}

func (receiver fakeEventbusOkHandler) Handle(ctx context.Context, payload string) error {
	select {
	case receiver.payloadChan <- []byte(payload):
	default:

	}
	return nil
}

func (receiver fakeEventbusOkHandler) IsOk() bool {
	select {
	case expected := <-receiver.payloadChan:
		return string(expected) == string(receiver.expectedPayload)
	case <-time.Tick(time.Second * 5):
		return false
	}
}

type fakeEventbusFailedHandler struct {
	typeA    chan []byte
	expected []byte
}

func (receiver fakeEventbusFailedHandler) Handle(ctx context.Context, payload string) error {
	select {
	case receiver.typeA <- []byte(""):
	default:
	}
	return fmt.Errorf("an error")
}

func (receiver fakeEventbusFailedHandler) IsOk() bool {
	select {
	case expected := <-receiver.typeA:
		return string(expected) == string(receiver.expected)
	case <-time.Tick(time.Second * 5):
		return false
	}
}
