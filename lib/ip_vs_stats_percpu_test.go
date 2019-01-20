package mpipvscpustat

import (
  "strings"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
  var IPVSCpustat IPVSCpustatPlugin

  graphdef := IPVSCpustat.GraphDefinition()
  if len(graphdef) != 3 {
    t.Errorf("GraphDefinition: %d should be 6", len(graphdef))
  }
}

func TestParse(t *testing.T) {
  var r1 IPVSCpustatPlugin
  r1.CPUs = 20
  s1 := `       Total Incoming Outgoing         Incoming         Outgoing
CPU    Conns  Packets  Packets            Bytes            Bytes
  0    EC76C   2A6BB9        0          AD77113                0
  1      C63     2127        0            87C64                0
  2      C77     21BC        0            8D081                0
  3      CED     2545        0            9A216                0
  4      C2F     1FC5        0            85794                0
  5      CDA     250B        0            98792                0
  6      C3E     22FD        0            8F7A1                0
  7      BE4     1FA0        0            85BE8                0
  8      CB7     2243        0            8C3CE                0
  9      CDD     239F        0            91325                0
  A      C48     216A        0            8A10D                0
  B      CB6     2314        0            91080                0
  C      CE8     2468        0            96958                0
  D      CB9     2325        0            90035                0
  E      C25     2299        0            8F942                0
  F      C3E     22E6        0            8E62B                0
 10      C23     225F        0            8C3FF                0
 11      D90     286C        0            A6358                0
 12      C1E     20A2        0            89FEF                0
 13      C84     244A        0            965E7                0
  ~    FB549   2D0271        0          B829164                0

     Conns/s   Pkts/s   Pkts/s          Bytes/s          Bytes/s
           0        0        0                0                0
`

  stubData1 := strings.NewReader(s1)
  a, err := r1.Parse(stubData1)
  assert.Nil(t, err)

  //        Total Incoming Outgoing         Incoming         Outgoing
  // CPU    Conns  Packets  Packets            Bytes            Bytes
  //   0    EC76C   2A6BB9        0          AD77113                0
  assert.EqualValues(t, 968556,    a["proc.net.ip_vs_stats_percpu.0.conns.conn"])
  assert.EqualValues(t, 2780089,   a["proc.net.ip_vs_stats_percpu.0.packets.in"])
  assert.EqualValues(t, 0,         a["proc.net.ip_vs_stats_percpu.0.packets.out"])
  assert.EqualValues(t, 181891347, a["proc.net.ip_vs_stats_percpu.0.bytes.in"])
  assert.EqualValues(t, 0,         a["proc.net.ip_vs_stats_percpu.0.bytes.out"])

  //        Total Incoming Outgoing         Incoming         Outgoing
  // CPU    Conns  Packets  Packets            Bytes            Bytes
  //   1      C63     2127        0            87C64                0
  assert.EqualValues(t, 3171,   a["proc.net.ip_vs_stats_percpu.1.conns.conn"])
  assert.EqualValues(t, 8487,   a["proc.net.ip_vs_stats_percpu.1.packets.in"])
  assert.EqualValues(t, 0,      a["proc.net.ip_vs_stats_percpu.1.packets.out"])
  assert.EqualValues(t, 556132, a["proc.net.ip_vs_stats_percpu.1.bytes.in"])
  assert.EqualValues(t, 0,      a["proc.net.ip_vs_stats_percpu.1.bytes.out"])

  //        Total Incoming Outgoing         Incoming         Outgoing
  // CPU    Conns  Packets  Packets            Bytes            Bytes
  //  13      C84     244A        0            965E7                0
  assert.EqualValues(t, 3204,   a["proc.net.ip_vs_stats_percpu.19.conns.conn"])
  assert.EqualValues(t, 9290,   a["proc.net.ip_vs_stats_percpu.19.packets.in"])
  assert.EqualValues(t, 0,      a["proc.net.ip_vs_stats_percpu.19.packets.out"])
  assert.EqualValues(t, 615911, a["proc.net.ip_vs_stats_percpu.19.bytes.in"])
  assert.EqualValues(t, 0,      a["proc.net.ip_vs_stats_percpu.19.bytes.out"])

  // 20 (LvsCpustat.CPUs) * 5 (fields)
  assert.EqualValues(t, 100, len(a))

  // CPU core is 2
  var r IPVSCpustatPlugin
  r.CPUs = 2
  s2 := `       Total Incoming Outgoing         Incoming         Outgoing
CPU    Conns  Packets  Packets            Bytes            Bytes
  0    EC76C   2A6BB9        0          AD77113                0
  1      C63     2127        0            87C64                0
  2      C77     21BC        0            8D081                0
  3      CED     2545        0            9A216                0
  4      C2F     1FC5        0            85794                0
  5      CDA     250B        0            98792                0
  6      C3E     22FD        0            8F7A1                0
  7      BE4     1FA0        0            85BE8                0
  8      CB7     2243        0            8C3CE                0
  9      CDD     239F        0            91325                0
  A      C48     216A        0            8A10D                0
  B      CB6     2314        0            91080                0
  C      CE8     2468        0            96958                0
  D      CB9     2325        0            90035                0
  E      C25     2299        0            8F942                0
  F      C3E     22E6        0            8E62B                0
 10      C23     225F        0            8C3FF                0
 11      D90     286C        0            A6358                0
 12      C1E     20A2        0            89FEF                0
 13      C84     244A        0            965E7                0
  ~    FB549   2D0271        0          B829164                0

     Conns/s   Pkts/s   Pkts/s          Bytes/s          Bytes/s
           0        0        0                0                0
`

  stubData2 := strings.NewReader(s2)

  b, err := r.Parse(stubData2)
  assert.Nil(t, err)
  assert.EqualValues(t, 10, len(b))
  //        Total Incoming Outgoing         Incoming         Outgoing
  // CPU    Conns  Packets  Packets            Bytes            Bytes
  //   0    EC76C   2A6BB9        0          AD77113                0
  assert.EqualValues(t, 968556,    b["proc.net.ip_vs_stats_percpu.0.conns.conn"])
  assert.EqualValues(t, 2780089,   b["proc.net.ip_vs_stats_percpu.0.packets.in"])
  assert.EqualValues(t, 0,         b["proc.net.ip_vs_stats_percpu.0.packets.out"])
  assert.EqualValues(t, 181891347, b["proc.net.ip_vs_stats_percpu.0.bytes.in"])
  assert.EqualValues(t, 0,         b["proc.net.ip_vs_stats_percpu.0.bytes.out"])

  //        Total Incoming Outgoing         Incoming         Outgoing
  // CPU    Conns  Packets  Packets            Bytes            Bytes
  //   1      C63     2127        0            87C64                0
  assert.EqualValues(t, 3171,   b["proc.net.ip_vs_stats_percpu.1.conns.conn"])
  assert.EqualValues(t, 8487,   b["proc.net.ip_vs_stats_percpu.1.packets.in"])
  assert.EqualValues(t, 0,      b["proc.net.ip_vs_stats_percpu.1.packets.out"])
  assert.EqualValues(t, 556132, b["proc.net.ip_vs_stats_percpu.1.bytes.in"])
  assert.EqualValues(t, 0,      b["proc.net.ip_vs_stats_percpu.1.bytes.out"])
}

func TestHex2Float64(t *testing.T) {
  a, err := Hex2Float64("0050")
  assert.Nil(t, err)
  assert.EqualValues(t, 80, a)
}

func TestCPUStatData(t *testing.T) {
  a, err := CPUStatData([]string{"F","C3E","22E6","0","8E62B","0",})
  assert.Nil(t, err)
  assert.EqualValues(t, 3134, a.Conns)
  assert.EqualValues(t, 8934, a.InPackets)
  assert.EqualValues(t, 0, a.OutPackets)
  assert.EqualValues(t, 583211, a.InBytes)
  assert.EqualValues(t, 0, a.OutBytes)
}
