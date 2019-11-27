// +build darwin freebsd netbsd

package network

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

// Get network statistics
func Get() ([]Stats, error) {
	// Reference: man 1 netstat
	cmd := exec.Command("netstat", "-bni")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	networks, err := collectNetworkStats(out)
	if err != nil {
		cmd.Process.Kill() // The process may stuck on write to pipe if the pipe buffer in kernel is full.
		go cmd.Wait()
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return networks, nil
}

// Stats represents network statistics for darwin
type Stats struct {
	Name             string
	RxBytes, TxBytes uint64
}

func collectNetworkStats(out io.Reader) ([]Stats, error) {
	scanner := bufio.NewScanner(out)

	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to scan output of netstat")
	}
	line := scanner.Text()
	if !strings.HasPrefix(line, "Name") {
		return nil, fmt.Errorf("unexpected output of netstat -bni: %s", line)
	}
	var rxBytesIdx, txBytesIdx int
	fields := strings.Fields(line)
	fieldsCount := len(fields)
	for i, field := range fields {
		switch field {
		case "Ibytes":
			rxBytesIdx = i
		case "Obytes":
			txBytesIdx = i
		}
	}
	if rxBytesIdx == 0 || txBytesIdx == 0 {
		return nil, fmt.Errorf("unexpected output of netstat -bni: %s", line)
	}

	var networks []Stats
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		name := strings.TrimSuffix(fields[0], "*")
		if strings.HasPrefix(name, "lo") || !strings.HasPrefix(fields[2], "<Link#") {
			continue
		}
		rxBytesIdx, txBytesIdx := rxBytesIdx, txBytesIdx
		if len(fields) < fieldsCount { // Address can be empty
			rxBytesIdx, txBytesIdx = rxBytesIdx-1, txBytesIdx-1
		}
		rxBytes, err := strconv.ParseUint(fields[rxBytesIdx], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Ibytes of %s", name)
		}
		txBytes, err := strconv.ParseUint(fields[txBytesIdx], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Obytes of %s", name)
		}
		networks = append(networks, Stats{Name: name, RxBytes: rxBytes, TxBytes: txBytes})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error for netstat: %s", err)
	}

	return networks, nil
}
