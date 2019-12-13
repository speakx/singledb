package genserver

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool"
)

type codecServer struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *pool.WorkerPool
}

func (cs *codecServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Test codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumLoops)
	return
}

func (cs *codecServer) React(c gnet.Conn) (out []byte, action gnet.Action) {
	if cs.async {
		data := append([]byte{}, c.ReadFrame()...)
		_ = cs.workerPool.Submit(func() {
			c.AsyncWrite(data)
		})
		return
	}
	out = c.ReadFrame()
	fmt.Println(out)
	return
}

func TestCodecServe(addr string, multicore, async bool, codec gnet.ICodec) {
	var err error
	if codec == nil {
		encoderConfig := gnet.EncoderConfig{
			ByteOrder:                       binary.BigEndian,
			LengthFieldLength:               4,
			LengthAdjustment:                0,
			LengthIncludesLengthFieldLength: false,
		}
		decoderConfig := gnet.DecoderConfig{
			ByteOrder:           binary.BigEndian,
			LengthFieldOffset:   0,
			LengthFieldLength:   4,
			LengthAdjustment:    0,
			InitialBytesToStrip: 4,
		}
		codec = gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)
	}
	cs := &codecServer{addr: addr, multicore: multicore, async: async, codec: codec, workerPool: pool.NewWorkerPool()}
	err = gnet.Serve(cs, addr, gnet.WithMulticore(multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec))
	if err != nil {
		panic(err)
	}
}
