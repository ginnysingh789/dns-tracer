# DNS Tracer

A simple command-line tool to perform an iterative DNS query and trace the resolution path from the root servers down to the final answer. This helps in understanding and debugging DNS resolution issues by showing each step of the process.

## Features

*   **Iterative Resolution:** Starts from the root servers and follows referrals.
*   **Trace Visualization:** Displays the query path in a clear, tree-like structure.
*   **CNAME Handling:** Follows CNAME records and restarts the resolution process.
*   **TCP Fallback:** Automatically switches to TCP for truncated UDP responses.
*   **Detailed Information:** Shows response times, server types (Root, TLD, Authoritative), and final IP addresses.

## How It Works

The tool mimics the behavior of a recursive DNS resolver when it needs to find an IP address for a domain it doesn't have cached. The process is as follows:

1.  **Query a Root Server:** It starts by asking one of the hardcoded root DNS servers for the domain's IP address.
2.  **Follow Referrals:** The root server won't know the final answer but will refer the tool to the appropriate Top-Level Domain (TLD) server (e.g., the `.com` or `.net` servers).
3.  **Query the TLD Server:** The tool then asks the TLD server. The TLD server, in turn, refers it to the authoritative name server for that specific domain. This is the server that holds the actual DNS records.
4.  **Query the Authoritative Server:** Finally, the tool queries the authoritative server, which provides the final answer—usually an IP address (A or AAAA record) or a CNAME record pointing to another domain.
5.  **Display the Path:** Each step of this journey is recorded and displayed, giving you a complete picture of the resolution path.

## Installation

To build the tool, you need to have Go installed on your system.

```sh
go build -o dns-tracer ./cmd
```

This will create an executable file named `dns-tracer` in the root directory.

## Usage

Run the tool from your terminal, providing the domain you want to trace using the `-domain` flag.

```sh
./dns-tracer -domain www.github.com
```

### Example Output

```
DNS Resolution Path
==================

ROOT: Asking 198.41.0.4 about www.github.com A
       (a.root-servers.net)
       Used TCP due to large response
       Response time: 477.9155ms
       Server doesn't know, but suggests asking:
         l.gtld-servers.net
         j.gtld-servers.net
         (and 11 others)
       Here are their IP addresses:
         l.gtld-servers.net at 192.41.162.30
         j.gtld-servers.net at 192.48.79.30
         (and 11 more)

  └─ TLD: Asking 192.41.162.30 about www.github.com A
         Response time: 159.4817ms
         Server doesn't know, but suggests asking:
           ns-520.awsdns-01.net
           ns-421.awsdns-52.com
           (and 6 others)
         Here are their IP addresses:
           ns-421.awsdns-52.com at 205.251.193.165

    └─ AUTH: Asking 205.251.193.165 about www.github.com A
           Response time: 7.8126ms
           Got answer:
             www.github.com is actually github.com
             Need to look up the real name now...

ROOT: Asking 198.41.0.4 about github.com A
       (a.root-servers.net)
       Used TCP due to large response
       Response time: 463.8984ms
       Server doesn't know, but suggests asking:
         l.gtld-servers.net
         j.gtld-servers.net
         (and 11 others)
       Here are their IP addresses:
         l.gtld-servers.net at 192.41.162.30
         j.gtld-servers.net at 192.48.79.30
         (and 11 more)

  └─ TLD: Asking 192.41.162.30 about github.com A
         Response time: 156.9456ms
         Server doesn't know, but suggests asking:
           ns-520.awsdns-01.net
           ns-421.awsdns-52.com
           (and 6 others)
         Here are their IP addresses:
           ns-421.awsdns-52.com at 205.251.193.165

    └─ AUTH: Asking 205.251.193.165 about github.com A
           Response time: 5.8067ms
           Got answer:
             Final IP address: 20.207.73.82

Summary
-------
Total queries: 6
Total time: 1.2718605s
Final result: 20.207.73.82
```

## Project Structure

The codebase is organized to separate concerns:

-   `cmd/`: Contains the main application entry point.
-   `pkg/dnsresolver/`: Implements the core DNS resolution logic.
-   `pkg/display/`: Handles the presentation of the trace results.
-   `pkg/util/`: Provides helper functions used across the project.
