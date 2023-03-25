package connection

import (
	"errors"
	"net"
	"strconv"

	"github.com/mortedecai/rivulets-of-go/logger"
)

// Data is the individual RoG data for a particular *net.Conn, including the connection itself.
type Data struct {
	net.Conn
}

// HandlerFunc operates against the connection.Data object and may modify it as needed.
type HandlerFunc func(data *Data) *Data

// Manager provides the methods required for a RoG Manager
// These include the abilty to Start, Start in Maintenance mode and Stop the listener
type Manager interface {
	// MaintenanceStart listens for incoming connections, responds with the maintenance message and closes the connection immediately.
	MaintenanceStart() error
	// SetHandler sets the default handler to invoke when not in maintenance mode.
	SetHandler(f HandlerFunc)
	// SetMaintenanceHandler sets the handler to use when maintenance mode is turned on.
	SetMaintenanceHandler(f HandlerFunc)
	// Start begins to listen for incoming connections
	Start() error
	// Stop prevents any further connections from being accepted.
	Stop()
}

type rogManager struct {
	address      *net.TCPAddr
	connections  []*Data
	listener     net.Listener
	logger       logger.Logger
	terminateMUD bool
	maintenance  bool

	handler            HandlerFunc
	maintenanceHandler HandlerFunc
}

// NewManager creates a new connection Manager with an identifying logger name and listening on the provided port
//     The port may be provided as an integer string (e.g. "3160") or in address specification (e.g. ":3160").
func NewManager(port string, l logger.Logger) (Manager, error) {
	var err error
	mgr := &rogManager{
		address:     new(net.TCPAddr),
		connections: make([]*Data, 0),
		logger:      l.WithName("mgr"),
	}

	mgr.logger.Infow("Setting Port", "Port", port)
	if port[0] == ':' {
		port = port[1:]
	}
	if (*mgr.address).Port, err = strconv.Atoi(port); err != nil {
		mgr.logger.Errorw("Creating Port", "Conversion error", err)
		return nil, err
	}

	return mgr, nil
}

func (rm *rogManager) Start() error {
	var err error
	rm.terminateMUD = false
	if rm.listener, err = net.ListenTCP("tcp", rm.address); err != nil {
		rm.logger.Errorw("Start Listening", "Listen error", err)
		rm.terminateMUD = true
		return err
	}

	go rm.listen()
	return nil
}

func (rm *rogManager) MaintenanceStart() error {
	rm.maintenance = true
	return rm.Start()
}

func (rm *rogManager) Stop() {
	rm.logger.Infow("Shutting Down")
	rm.terminateMUD = true
	rm.listener.Close()
}

func (rm *rogManager) SetHandler(f HandlerFunc) {
	rm.handler = f
}

func (rm *rogManager) SetMaintenanceHandler(f HandlerFunc) {
	rm.maintenanceHandler = f
}

func (rm *rogManager) listen() {
	const (
		methodName  = "listen()"
		maintenance = "Maintenance"
	)
	rm.logger.Infow(methodName, maintenance, rm.maintenance)
	var conn net.Conn
	var err error
	for {
		if rm.terminateMUD {
			break
		}
		conn, err = rm.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			rm.logger.Infow(methodName, "Listener Status", "Closed")
			break
		} else if err != nil {
			rm.logger.Errorw(methodName, "Accept Error", err)
			continue
		}

		if rm.maintenance {
			rm.maintenanceHandler(&Data{conn})
			rm.logger.Debugw(methodName, maintenance, rm.maintenance, "Remote Address", conn.(*net.TCPConn).RemoteAddr())
			if err = conn.Close(); err != nil {
				rm.logger.Errorw("Close Connection", "Error", err)
			}
		} else {
			data := &Data{conn}
			rm.connections = append(rm.connections, data)
			rm.logger.Debugw(methodName, maintenance, rm.maintenance, "Count", len(rm.connections), "Remote Address", conn.(*net.TCPConn).RemoteAddr())
			rm.handler(data)
		}
	}
}
