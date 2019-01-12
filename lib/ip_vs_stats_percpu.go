package mpipvscpustat

import(
  "strings"
  "flag"
  "io"
  "bufio"
  "os"
  "strconv"
  "fmt"
  "runtime"

  mp "github.com/mackerelio/go-mackerel-plugin"
)

// IPVSCpustatPlugin struct
type IPVSCpustatPlugin struct {
  Prefix string
  CPUs int
}

// IPVSCpustatData struct
type IPVSCpustatData struct {
  Conns float64
  InPackets float64
  OutPackets float64
  InBytes float64
  OutBytes float64
}

// KeyPrefixTemplate ...
var KeyPrefixTemplate = "proc.net.ip_vs_stats_percpu.#"

// define graph
var graphdef = map[string]mp.Graphs{
  KeyPrefixTemplate + ".conns": {
    Label: "ip_vs_stats_percpu Connections",
    Unit: "integer",
    Metrics: []mp.Metrics {
      {Name: "conn", Label: "Connections", Diff: true },
    },
  },
  KeyPrefixTemplate + ".packets": {
    Label: "ip_vs_stats_percpu Packets",
    Unit: "integer",
    Metrics: []mp.Metrics {
      { Name: "in", Label: "Incoming Packets", Diff: true },
      { Name: "out", Label: "Outgoing Packets", Diff: true },
    },
  },
  KeyPrefixTemplate + ".bytes": {
    Label: "ip_vs_stats_percpu Bytes",
    Unit: "bytes",
    Metrics: []mp.Metrics {
      { Name: "in", Label: "Incoming Packets", Diff: true },
      { Name: "out", Label: "Outgoing Packets", Diff: true },
    },
  },
}

// GraphDefinition : interface for go-mackerel-plugin interface
func (r IPVSCpustatPlugin) GraphDefinition() map[string]mp.Graphs {
  return graphdef
}

// FetchMetrics : interface for go-mackerel-plugin interface
func (r IPVSCpustatPlugin) FetchMetrics() (map[string]float64, error) {
  file, err := os.Open("/proc/net/ip_vs_stats_percpu")
  if err != nil {
    return nil, err
  }
  defer file.Close()

  return r.Parse(file)
}

// Parse : parser for /proc/net/ip_vs_stats_percpu for mackerel-plugin-lvs-cpustat
func (r IPVSCpustatPlugin) Parse(stat io.Reader) (map[string]float64, error) {
  scanner := bufio.NewScanner(stat)
  data := make(map[string]float64)

  for scanner.Scan() {
    fields := strings.Fields(scanner.Text())
    // cpu stat line is follow format
    // <CpuNum> <Conns> <InPackets> <OutPackets> <InBytes> <OutBytes>
    if len(fields) != 6 {
      continue
    }
    if fields[0] == "CPU" || fields[0] == "~" {
      continue
    }
    CPUNum, err := strconv.ParseUint(fields[0], 16, 64)
    if err != nil {
      return nil, err
    }
    if int(CPUNum) >= r.CPUs {
      continue
    }
    KeyPrefix := strings.Replace(KeyPrefixTemplate, "#", fmt.Sprint(CPUNum), -1)
    d, err := CPUStatData(fields)
    if err != nil {
      return nil, err
    }
    data[KeyPrefix + ".conns.conn"] = d.Conns
    data[KeyPrefix + ".packets.in"] = d.InPackets
    data[KeyPrefix + ".packets.out"] = d.OutPackets
    data[KeyPrefix + ".bytes.in"] = d.InBytes
    data[KeyPrefix + ".bytes.out"] = d.OutBytes
  }

  return data, nil
}

// CPUStatData ...
func CPUStatData(fields []string) (IPVSCpustatData, error) {
  var (
    d IPVSCpustatData
    err error
  )

  d.Conns, err = Hex2Float64(fields[1])
  if err != nil {
    return d, err
  }
  d.InPackets, err = Hex2Float64(fields[2])
  if err != nil {
    return d, err
  }
  d.OutPackets, err = Hex2Float64(fields[3])
  if err != nil {
    return d, err
  }
  d.InBytes, err = Hex2Float64(fields[4])
  if err != nil {
    return d, err
  }
  d.OutBytes, err = Hex2Float64(fields[5])
  if err != nil {
    return d, err
  }

  return d, nil
}

// Hex2Float64 : convert hex string to float64
func Hex2Float64(hex string) (float64, error) {
  a, err := strconv.ParseUint(hex, 16, 64)
  if err != nil {
    return 0.0, err
  }
  return float64(a), nil
}

// Do the plugin
func Do() {
  optCPUs := flag.Int("cpus", 0, "Count of CPU cores")
  optTempfile := flag.String("tempfile", "", "Temp file name")
  flag.Parse()

  var IPVSCpustat IPVSCpustatPlugin
  if *optCPUs == 0 {
    IPVSCpustat.CPUs = runtime.NumCPU()
  } else {
    IPVSCpustat.CPUs = *optCPUs
  }

  helper := mp.NewMackerelPlugin(IPVSCpustat)
  helper.Tempfile = *optTempfile

  helper.Run()
}
