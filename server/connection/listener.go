package connection

import (
	"net"
	"strconv"

	"go.uber.org/zap"
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
	MaintenanceStart()
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
	listener     *net.TCPListener
	logger       *zap.SugaredLogger
	terminateMUD bool
	maintenance  bool

	handler            HandlerFunc
	maintenanceHandler HandlerFunc
}

// NewManager creates a new connection Manager with an identifying logger name and listening on the provided port
//     The port may be provided as an integer string (e.g. "3160") or in address specification (e.g. ":3160").
func NewManager(port string, logger *zap.SugaredLogger) (Manager, error) {
	var err error
	mgr := &rogManager{
		address:     new(net.TCPAddr),
		connections: make([]*Data, 0),
		logger:      logger.Named("rogManager"),
	}

	mgr.logger.Infow("Creating Port", "Port", port)
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
		rm.logger.Errorw("Creating Port", "Listen error", err)
		rm.terminateMUD = true
		return err
	}

	go rm.listen()
	return nil
}

func (rm *rogManager) MaintenanceStart() {
	rm.maintenance = true
	rm.Start()
}

func (rm *rogManager) Stop() {
	rm.logger.Infow("Shutting Down")
	rm.terminateMUD = true
}

func (rm *rogManager) SetHandler(f HandlerFunc) {
	rm.handler = f
}

func (rm *rogManager) SetMaintenanceHandler(f HandlerFunc) {
	rm.maintenanceHandler = f
}

func (rm *rogManager) listen() {
	rm.logger.Infow("Listening for incoming connections")
	var conn net.Conn
	var err error
	for {
		if conn, err = rm.listener.Accept(); err != nil {
			rm.logger.Errorw("Error accepting connection", "Error", err)
			continue
		}
		if rm.maintenance {
			rm.maintenanceHandler(&Data{conn})
			conn.Close()
			rm.logger.Debugw("Accepted Connection", "Mode", "maintenance", "Remote Address", conn.(*net.TCPConn).RemoteAddr())
		} else {
			rm.connections = append(rm.connections, &Data{conn})
			rm.logger.Debugw("Accepted Connection", "Mode", "regular", "Count", len(rm.connections), "Remote Address", conn.(*net.TCPConn).RemoteAddr())
		}
		if rm.terminateMUD {
			break
		}
	}
}
