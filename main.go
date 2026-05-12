package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
    "time"

    "github.com/jackpal/go-nat-pmp"
    "github.com/urfave/cli/v3"
    "golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
)

func main() {
    cmd := &cli.Command{
        Name:        "natpmpc",
        Usage:       "NAT-PMP keep-alive tool",
        Description: "Continuously refreshes UDP+TCP port mappings every 45 seconds using NAT-PMP.",
        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:    "help",
                Aliases: []string{"h"},
                Usage:   "display this help screen",
            },
            &cli.StringFlag{
                Name:    "gateway",
                Aliases: []string{"g"},
                Usage:   "force the gateway IPv4 address to use",
            },
        },
        Action: func(ctx context.Context, cmd *cli.Command) error {
            if cmd.Bool("help") {
                cli.ShowAppHelp(cmd)
                return nil
            }
            return runKeepAliveLoop(cmd)
        },
    }

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}

func getWireGuardGateway() (net.IP, error) {
    fmt.Println("Scanning network adapters for WireGuard Tunnel...")

    adapters, err := winipcfg.GetAdaptersAddresses(0, 0)
    if err != nil {
        return nil, fmt.Errorf("failed to get adapters: %w", err)
    }

    for _, adapter := range adapters {
        if strings.EqualFold(adapter.Description(), "WireGuard Tunnel") {
            fmt.Println("Found WireGuard Tunnel interface")
            // Try gateway address first
            gw := adapter.FirstGatewayAddress
            for gw != nil {
                if ip := gw.Address.IP(); ip.To4() != nil {
                    return ip, nil
                }
                gw = gw.Next
            }

            fmt.Println("Gateway IP unknown, using derivation method...")

            // Fallback: derive .1 from interface IP (common on WireGuard)
            ua := adapter.FirstUnicastAddress
            for ua != nil {
                if ip := ua.Address.IP(); ip.To4() != nil {
                    if len(ip) == 4 {
                        gwIP := make(net.IP, len(ip))
                        copy(gwIP, ip)
                        gwIP[3] = 1 // change last octet to .1
                        return gwIP, nil
                    }
                }
                ua = ua.Next
            }

            return nil, fmt.Errorf("WireGuard Tunnel found but no usable IPv4 address")
        }
    }

    return nil, fmt.Errorf("WireGuard Tunnel interface not found")
}

func getClient(gwFlag string) (*natpmp.Client, error) {
    var gwIP net.IP
    var err error

    if gwFlag != "" {
        gwIP = net.ParseIP(gwFlag)
        if gwIP == nil || gwIP.To4() == nil {
            return nil, fmt.Errorf("invalid gateway IP: %s", gwFlag)
        }
    } else {
        gwIP, err = getWireGuardGateway()
        if err != nil {
            return nil, err
        }
    }

    fmt.Printf("Using gateway IP %s\n\n", gwIP)
    return natpmp.NewClient(gwIP), nil
}

func runKeepAliveLoop(cmd *cli.Command) error {
    gwFlag := cmd.String("gateway")

    fmt.Println("Starting NAT-PMP keep-alive loop; refresh UDP+TCP port every 45s forever")
    fmt.Println("Press Ctrl+C to stop.\n")

    client, err := getClient(gwFlag)
    if err != nil {
        return err
    }

    for {
        fmt.Println(time.Now().Format(time.RFC1123))

        if err := doMapping(client, 1, 0, "udp", 60); err != nil {
            fmt.Printf("ERROR with UDP mapping: %v\n", err)
            break
        }

        if err := doMapping(client, 1, 0, "tcp", 60); err != nil {
            fmt.Printf("ERROR with TCP mapping: %v\n", err)
            break
        }

        time.Sleep(45 * time.Second)
    }
    return nil
}

func doMapping(client *natpmp.Client, public, private int, proto string, lifetime int) error {
    result, err := client.AddPortMapping(proto, private, public, lifetime)
    if err != nil {
        return err
    }

    protoName := strings.ToUpper(proto)
    fmt.Printf("Mapped public port %d protocol %s to local port %d lifetime %d\n",
        result.MappedExternalPort, protoName, result.InternalPort, result.PortMappingLifetimeInSeconds)

    return nil
}
