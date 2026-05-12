# go-natpmpc

## Real-world example

The below is using WireGuard with ProtonVPN:

```
D:\>go-natpmpc.exe
Starting NAT-PMP keep-alive loop; refresh UDP+TCP port every 45s forever
Press Ctrl+C to stop.

Scanning network adapters for WireGuard Tunnel...
Found WireGuard Tunnel interface
Gateway IP unknown, using derivation method...
Using gateway IP 10.2.0.1

Tue, 12 May 2026 10:04:28 PDT
Mapped public port 59631 protocol UDP to local port 0 lifetime 60
Mapped public port 59631 protocol TCP to local port 0 lifetime 60
Tue, 12 May 2026 10:05:13 PDT
Mapped public port 59631 protocol UDP to local port 0 lifetime 60
Mapped public port 59631 protocol TCP to local port 0 lifetime 60
Tue, 12 May 2026 10:05:58 PDT
Mapped public port 59631 protocol UDP to local port 0 lifetime 60
Mapped public port 59631 protocol TCP to local port 0 lifetime 60
...
```

## Usage

```
D:\>go-natpmpc.exe -h
NAME:
   natpmpc - NAT-PMP keep-alive tool

USAGE:
   natpmpc [global options]

DESCRIPTION:
   Continuously refreshes UDP+TCP port mappings every 45 seconds using NAT-PMP.

GLOBAL OPTIONS:
   --help, -h                   display this help screen
   --gateway string, -g string  force the gateway IPv4 address to use
   --help, -h                   show help
```

## Installing from source on Windows

```
go install github.com/koitsu/go-natpmpc@v1.0.0
%GOPATH%\bin\go-natpmpc.exe
```

Linux and OS X should be identical, though replace `go-natpmpc.exe` with `go-natpmpc`.

# Why I did this

My use-case was as follows:

* Windows 10 Pro as the OS
* WireGuard for the VPN software (because ProtonVPN's Windows client is bloated, buggy, and generally awful; WireGuard in contrast is tiny and fantastic)
* ProtonVPN as the VPN provider -- which, requires use of NAT-PMP to ask for, and keep open, a port forward
* A single standalone binary that took care of the NAT-PMP aspect, as well as automatically determining the gateway IP address to use (if possible)

I was disgusted by the fact that:

1. There was no solution for Windows
1. [ProtonVPN's own instructions](https://protonvpn.com/support/port-forwarding-manual-setup#how-to-use-port-forwarding) only covered Linux and OS X, with different requirements for both
1. Not all Linux distros offer a libnatpmp package
1. Forcing reliance on Python on OS X for no good reason
1. Required use of sh/bash/a shell to run a `while true` loop to keep the port forward open.  (The OS X incarnation is even worse.)

# Disclaimers

**AI was used to create this software.**

I originally created this through the help of Grok.  I asked it to port [libnatpmp/natpmpc.c](https://github.com/miniupnp/libnatpmp/blob/master/natpmpc.c) from the [miniupnp project](https://github.com/miniupnp) into Go, using [jackpal/go-nat-pmp](https://github.com/jackpal/go-nat-pmp) as the NAT-PMP library.  I then made manual tweaks/QoL adjustments to get what I wanted.

**I am not intimately familiar with the Go language.**

While I'm old, know a lot of PLs (65xxx assembly, PIC16C84 assembly, x86 assembly (286/386/486 only), Pascal, Perl, C, PHP, Python), and am a UNIX systems administrator of over 30 years, Go is very new to me.  I was not interested in learning all the esoteric details of Go just to achieve this basic goal.  So I felt this was a reasonable use of AI.

**I am not intimately familiar with the Rust language.**

I did try using Rust for this instead of Go.  I had no desire to install Visual Studio Community Edition (1-2GBytes) just to get access to MSVC shims.  And while I did make a version that did not use MSVC (via installing the Rust package under MSYS2), its reliance on libunwind.dll (dynamically) was a turn-off.  (I found out days later that there was a way to statically include libunwind, but by then I had already gone back to Go.)  And, like Go, I was not interested in learning all the esoteric details of Rust just to achieve this basic goal.

**I care about software quality tremendously, but not in this case.**

* The code is probably janky and ridiculous.  It's AI code.  I'm OK with that in this case
* I have no interest in dealing with anything relating to IPv6
* I am not interested in providing "support" for this, so any support requests will be ignored and closed,
* If someone wants to submit PRs for this, fix bugs, etc. -- be my guest (or heck, maybe I'll just give you access to the repo!)
* The only focus is that the end result be a standalone binary that handles NAT-PMP port forwarding requests/refreshes

# My hopes

* That Proton AG themselves start recommending this software (for Linux and OS X, and as a Windows alternative for people who don't want to use their VPN software), instead of telling people to build Python junk or install libnatpmp just for what should be basic functionality
* That someone forks this repository and provide actual binaries for multiple OSes (Linux, OS X, and Windows) so that end-users can get on with their lives
* That someone can get the binary down to a smaller size than 6MBytes :)
