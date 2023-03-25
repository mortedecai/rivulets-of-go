package connection_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mortedecai/rivulets-of-go/logger"
	ml "github.com/mortedecai/rivulets-of-go/logger/mock_logger"
	"github.com/mortedecai/rivulets-of-go/server/connection"
	"github.com/mortedecai/rivulets-of-go/server/info"
)

const (
	port = ":3160"
)

var _ = Describe("Manager", func() {
	var _ = Describe("Initialization", func() {
		var (
			l *ml.MockLogger
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			l = ml.NewMockLogger(ctrl)
		})

		const expName = "rog"
		It("should return a non-nil manager properly initialized", func() {
			l.EXPECT().WithName(gomock.AssignableToTypeOf(expName)).Times(1).DoAndReturn(func(_ string) logger.Logger { return l })
			l.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Errorw(gomock.Any(), gomock.Any()).Times(0)
			mgr, err := connection.NewManager(port, l)
			Expect(err).ToNot(HaveOccurred())
			Expect(mgr).ToNot(BeNil())

		})
	})
	var _ = Describe("Running", Ordered, func() {
		var (
			l   *ml.MockLogger
			mgr connection.Manager
			err error
		)
		const (
			mudName = "rog"
		)
		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			l = ml.NewMockLogger(ctrl)
			l.EXPECT().WithName(gomock.AssignableToTypeOf(mudName)).AnyTimes().DoAndReturn(func(_ string) logger.Logger { return l })
			l.EXPECT().Debugw(gomock.Any(), gomock.Any()).Times(1)
			l.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Errorw(gomock.Any(), gomock.Any()).Times(0)
			mgr, err = connection.NewManager(port, l)
			mgr.SetHandler(testHandler)
			mgr.SetMaintenanceHandler(testMaintenanceHandler)
			Expect(err).ToNot(HaveOccurred())
		})
		It("should be possible to run in non-maintenance mode connect and stop the server", func() {
			Expect(mgr.Start()).ToNot(HaveOccurred())

			conn, err := net.Dial("tcp", "127.0.0.1:3160")
			Expect(err).ToNot(HaveOccurred())

			connbuf := bufio.NewReader(conn)

			str, err := connbuf.ReadString('\n')
			Expect(err).ToNot(HaveOccurred())
			Expect(str).To(Equal(handlerString))

			//Expect(isConnectionClosed(conn)).To(BeFalse())

			mgr.Stop()
		})
		It("should be possible to run in maintenance mode connect and stop the server", func() {
			Expect(mgr.MaintenanceStart()).ToNot(HaveOccurred())

			conn, err := net.Dial("tcp", "127.0.0.1:3160")
			Expect(err).ToNot(HaveOccurred())
			defer conn.Close()

			connbuf := bufio.NewReader(conn)

			str, err := connbuf.ReadString('\n')
			Expect(err).ToNot(HaveOccurred())
			Expect(str).To(Equal(maintenanceString))

			Expect(isConnectionClosed(conn)).To(BeTrue())

			mgr.Stop()
		})

	})
})

func isConnectionClosed(conn net.Conn) bool {
	oneByte := make([]byte, 1)
	_, err := conn.Read(oneByte)
	return errors.Is(err, io.EOF)
}

func writeMessage(conn *connection.Data, dataStr string) *connection.Data {
	data := []byte(dataStr)

	totalBytes := 0
	bw := 0

	for (totalBytes + bw) < len(data) {
		bw, err := conn.Write(data[totalBytes:])
		totalBytes += bw
		if err != nil {
			fmt.Println("Error writing hello string:  ", err.Error())
		}
	}
	return conn

}

var (
	handlerString     = fmt.Sprintf("%s - %s - %s\n", info.Name, info.Version, info.Commit)
	maintenanceString = fmt.Sprintf("%s - %s - %s MAINTENANCE\n", info.Name, info.Version, info.Commit)
)

func testHandler(conn *connection.Data) *connection.Data {
	return writeMessage(conn, handlerString)
}

func testMaintenanceHandler(conn *connection.Data) *connection.Data {
	return writeMessage(conn, maintenanceString)
}
