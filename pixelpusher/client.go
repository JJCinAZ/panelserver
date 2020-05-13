package pixelpusher

import (
	"fmt"
	"github.com/jjcinaz/panelserver/bytesbuffer"
	"image"
	"image/color"
	"net"
	"strings"
)

type Client struct {
	sequence   uint32
	network    string
	address    string
	rows, cols int
	conn       *net.UDPConn
}

func NewClient(network, address string, rows, cols int) (*Client, error) {
	if rows < 1 || rows > 1024 {
		return nil, fmt.Errorf("rows out of range, must be 1...1024")
	}
	if cols < 1 || cols > 1024 {
		return nil, fmt.Errorf("cols out of range, must be 1...1024")
	}
	c := new(Client)
	c.sequence = 1
	c.rows, c.cols = rows, cols
	c.network, c.address = network, address
	return c, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *Client) connect() error {
	if c.conn == nil {
		var (
			err  error
			addr *net.UDPAddr
		)
		if strings.IndexRune(c.address, ':') == -1 {
			c.address += ":5078"
		}
		if addr, err = net.ResolveUDPAddr(c.network, c.address); err != nil {
			return err
		}
		c.conn, err = net.DialUDP("udp", nil, addr)
		return err
	}
	return nil
}

func (c *Client) SendImage(img image.Image) error {
	err := c.connect()
	if err != nil {
		return err
	}
	buf := c.createBuffer(img)
	c.conn.SetWriteBuffer(buf.Size())
	_, err = c.conn.Write(buf.Bytes())
	return err
}

func (c *Client) createBuffer(img image.Image) bytesbuffer.Buffer {
	s := img.Bounds()
	buf, _ := bytesbuffer.NewBuffer(bytesbuffer.LittleEndian)
	// Grow the buffer to max to avoid multiple memory allocation calls
	rows := s.Max.Y - s.Min.Y
	buf.Grow((rows)*(s.Max.X-s.Min.X)*3 + rows + 4)
	buf.PutUint32(c.sequence)
	c.sequence++
	for y := s.Min.Y; y < s.Max.Y; y++ {
		buf.PutByte(byte(y))
		for x := s.Min.X; x < s.Max.X; x++ {
			p := img.At(x, y)
			buf.PutByte(p.(color.RGBA).R)
			buf.PutByte(p.(color.RGBA).G)
			buf.PutByte(p.(color.RGBA).B)
		}
	}
	return buf
}
