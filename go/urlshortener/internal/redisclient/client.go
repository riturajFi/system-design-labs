package redisclient

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	addr string
}

func New(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Get(ctx context.Context, key string) (string, bool, error) {
	conn, reader, err := c.open(ctx)
	if err != nil {
		return "", false, err
	}
	defer conn.Close()

	if err := writeCommand(conn, "GET", key); err != nil {
		return "", false, err
	}

	return readBulkString(reader)
}

func (c *Client) SetNX(ctx context.Context, key string, value string) (bool, error) {
	conn, reader, err := c.open(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	if err := writeCommand(conn, "SETNX", key, value); err != nil {
		return false, err
	}

	n, err := readInteger(reader)
	if err != nil {
		return false, err
	}

	return n == 1, nil
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	conn, reader, err := c.open(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := writeCommand(conn, "SET", key, value); err != nil {
		return err
	}

	return readSimpleString(reader)
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	conn, reader, err := c.open(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	if err := writeCommand(conn, "INCR", key); err != nil {
		return 0, err
	}

	return readInteger(reader)
}

func (c *Client) Del(ctx context.Context, key string) error {
	conn, reader, err := c.open(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := writeCommand(conn, "DEL", key); err != nil {
		return err
	}

	_, err = readInteger(reader)
	return err
}

func (c *Client) open(ctx context.Context) (net.Conn, *bufio.Reader, error) {
	dialer := &net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return nil, nil, fmt.Errorf("dial redis %s: %w", c.addr, err)
	}

	return conn, bufio.NewReader(conn), nil
}

func writeCommand(conn net.Conn, args ...string) error {
	var builder strings.Builder

	builder.WriteString("*")
	builder.WriteString(strconv.Itoa(len(args)))
	builder.WriteString("\r\n")

	for _, arg := range args {
		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(len(arg)))
		builder.WriteString("\r\n")
		builder.WriteString(arg)
		builder.WriteString("\r\n")
	}

	if _, err := conn.Write([]byte(builder.String())); err != nil {
		return fmt.Errorf("write redis command: %w", err)
	}

	return nil
}

func readBulkString(reader *bufio.Reader) (string, bool, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return "", false, fmt.Errorf("read redis response: %w", err)
	}

	switch prefix {
	case '$':
		line, err := readLine(reader)
		if err != nil {
			return "", false, err
		}

		size, err := strconv.Atoi(line)
		if err != nil {
			return "", false, fmt.Errorf("parse bulk string size %q: %w", line, err)
		}
		if size == -1 {
			return "", false, nil
		}

		buf := make([]byte, size+2)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return "", false, fmt.Errorf("read bulk string payload: %w", err)
		}

		return string(buf[:size]), true, nil
	case '-':
		line, err := readLine(reader)
		if err != nil {
			return "", false, err
		}
		return "", false, fmt.Errorf("redis error: %s", line)
	default:
		return "", false, fmt.Errorf("unexpected redis response prefix: %q", prefix)
	}
}

func readInteger(reader *bufio.Reader) (int64, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read redis response: %w", err)
	}

	switch prefix {
	case ':':
		line, err := readLine(reader)
		if err != nil {
			return 0, err
		}

		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse integer response %q: %w", line, err)
		}

		return n, nil
	case '-':
		line, err := readLine(reader)
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("redis error: %s", line)
	default:
		return 0, fmt.Errorf("unexpected redis response prefix: %q", prefix)
	}
}

func readSimpleString(reader *bufio.Reader) error {
	prefix, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read redis response: %w", err)
	}

	switch prefix {
	case '+':
		_, err := readLine(reader)
		return err
	case '-':
		line, err := readLine(reader)
		if err != nil {
			return err
		}
		return fmt.Errorf("redis error: %s", line)
	default:
		return fmt.Errorf("unexpected redis response prefix: %q", prefix)
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read redis line: %w", err)
	}

	return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"), nil
}
