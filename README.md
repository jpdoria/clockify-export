# About

This is a simple tool to get the time entries from Clockify, do a USD to PHP conversion, and calculate the total amount using the hourly rate.

# Usage

1. Download the binary from the release page.
1. Set the environment variables `CLOCKIFY_API_KEY` and `HOURLY_RATE_USD`.
1. Execute the binary.

```bash
export CLOCKIFY_API_KEY="foo...bar"
export HOURLY_RATE_USD="10"
./clockify-export-arm64-0.3.0
```

# Sample Output

```bash
→ ./clockify-export-arm64
Work Log:
ID     DATE           HOURS                EARNINGS
1      2024-05-01     HH:MM:SS (XX.YY)     $10.00
2      2024-05-02     HH:MM:SS (XX.YY)     $10.00
3      2024-05-03     HH:MM:SS (XX.YY)     $10.00
4      2024-05-06     HH:MM:SS (XX.YY)     $10.00
5      2024-05-07     HH:MM:SS (XX.YY)     $10.00
6      2024-05-08     HH:MM:SS (XX.YY)     $10.00
7      2024-05-09     HH:MM:SS (XX.YY)     $10.00
8      2024-05-10     HH:MM:SS (XX.YY)     $10.00
9      2024-05-13     HH:MM:SS (XX.YY)     $10.00
10     2024-05-14     HH:MM:SS (XX.YY)     $10.00
11     2024-05-15     HH:MM:SS (XX.YY)     $10.00
12     2024-05-16     HH:MM:SS (XX.YY)     $10.00
13     2024-05-17     HH:MM:SS (XX.YY)     $10.00
14     2024-05-20     HH:MM:SS (XX.YY)     $10.00
15     2024-05-21     HH:MM:SS (XX.YY)     $10.00
16     2024-05-22     HH:MM:SS (XX.YY)     $10.00
17     2024-05-23     HH:MM:SS (XX.YY)     $10.00
Total Hours: HH:MM:SS (XX.YY)

Exchange rate right now for 1 USD to PHP: 58.18
Total Earnings: $170.00 (P9890.60)
```

# TODO

- Prettier output. Export as PDF, instead of spreadsheet.
