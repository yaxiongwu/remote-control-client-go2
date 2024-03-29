package engine

import (
	"context"
	//"crypto/tls"
	//"crypto/x509"
	//"io/ioutil"

	"github.com/yaxiongwu/remote-control-client-go2/pkg/proto/rtc"

	log "github.com/pion/ion-log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ConnectorConfig struct {
	SSL    bool
	Cafile string
	Token  string
}

type ServiceEvent struct {
	Name      string
	ErrStatus *status.Status
}

type Service interface {
	Name() string
	Connect()
	Close()
	Connected() bool
}

type Connector struct {
	services map[string]Service
	Metadata metadata.MD

	config   *ConnectorConfig
	grpcConn *grpc.ClientConn

	OnOpen  func(Service)
	OnClose func(Service, ServiceEvent)

	ctx context.Context
}

// NewConnector create a ion connector
func NewConnector(addr string, config ...ConnectorConfig) *Connector {
	c := &Connector{
		services: make(map[string]Service),
		Metadata: make(metadata.MD),
		ctx:      context.Background(),
		config: &ConnectorConfig{
			SSL: true,
			//Cafile:"bxzryd.pem",
			Cafile: "../cred/bxzryd.cn_bundle.crt",
			//Cafile:"bxzryd.cn_bundle.pem",
			//Cafile:"bxzryd.cn_bundle.crt",
			//Token:"",
		},
	}

	if len(config) > 0 {
		c.config = &config[0]
	}

	//c.config.SSL=true
	//c.config.Cafile="bxzryd.pem"

	if addr == "" {
		log.Errorf("error: %v", errInvalidAddr)
		return nil
	}

	if c.config != nil && c.config.Token != "" {
		c.Metadata.Append("authorization", c.config.Token)
	}

	var err error
	if c.config != nil && c.config.SSL {
		/*
			var config *tls.Config
			if c.config.Cafile != "" {
				b, _ := ioutil.ReadFile(c.config.Cafile)
				cp := x509.NewCertPool()
				if !cp.AppendCertsFromPEM(b) {
					log.Errorf("credentials: failed to append certificates")
					return nil
			    }

				config = &tls.Config{
					InsecureSkipVerify: false,
					RootCAs:            cp,
				}
			} else {
				config = &tls.Config{
					InsecureSkipVerify: false,
				}
			}
			* */

		// creds,err:=credentials.NewClientTLSFromFile("bxzryd.cn_bundle.crt","")
		//creds,err:=credentials.NewClientTLSFromFile("bxzryd.cn_bundle.pem","")
		// creds,err:=credentials.NewClientTLSFromFile("bxzryd.cn.csr","")
		// creds,err:=credentials.NewClientTLSFromFile("bxzryd.cn.key","")
		/*
			      if err != nil {
						log.Errorf("credentials err: %v", err)
						return nil
					}

		*/

		creds, err := credentials.NewClientTLSFromFile("../cred/bxzryd.cn_bundle.crt", "")
		if err != nil {
			log.Errorf("NewClientTLSFromFile err: %v", err)
		}
		c.grpcConn, err = grpc.Dial("www.bxzryd.cn:5551", grpc.WithTransportCredentials(creds))

		//c.grpcConn, err = grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithBlock())
		//c.grpcConn, err = grpc.Dial(addr, grpc.WithTransportCredentials(creds), grpc.WithBlock())
		if err != nil {
			log.Errorf("did not connect: %v", err)
			return nil
		}

	} else {
		c.grpcConn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	}

	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil
	}

	log.Infof("gRPC connected: %s", addr)

	return c
}

func (c *Connector) Signal(r *RTC) (Signaller, error) {
	c.RegisterService(r)
	client := rtc.NewRTCClient(c.grpcConn)
	r.ctx = metadata.NewOutgoingContext(r.ctx, c.Metadata)
	return client.Signal(r.ctx)
}

func (c *Connector) Close() {
	for _, s := range c.services {
		if s.Connected() {
			s.Close()
		}
	}

}

func (c *Connector) OnHeaders(service Service, headers metadata.MD) {
	for k, v := range headers {
		c.Metadata.Append(k, v[0])
	}

	if c.OnOpen != nil {
		c.OnOpen(service)
	}
}

func (c *Connector) OnEnd(service Service, headers metadata.MD) {
	if c.OnClose != nil {
		c.OnClose(service, ServiceEvent{
			Name:      service.Name(),
			ErrStatus: status.New(codes.OK, "close"),
		})
	}
}

func (c *Connector) RegisterService(service Service) {
	c.services[service.Name()] = service
}
