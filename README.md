# About

This is a simple tool to get the time entries from Clockify, do a USD to PHP conversion, and calculate the total amount using the hourly rate.

# Usage

If you're not using Apple Silicon, please run the `clockify-export-amd64-*` executable instead.

1. Download the binary from the release page.
1. Set the environment variables `CLOCKIFY_API_KEY` and `HOURLY_RATE_USD`.
1. Execute the binary.

```bash
export CLOCKIFY_API_KEY="foo...bar"
export HOURLY_RATE_USD="10"
./clockify-export-arm64-0.4.4
```

## Without `-customRange` flag (defaults to current month)

```bash
❯ ./clockify-export-arm64-0.4.4
```

## With `-customRange` flag

```bash
❯ ./clockify-export-arm64-0.4.4 -customRange "2024-05-01T00:00:00.000Z to 2024-05-31T23:59:59.999Z"
```

# Sample Output

```bash
❯ ./clockify-export-arm64-0.4.4
Work Log:
ID     DATE           HOURS                EARNINGS
1      2024-05-01     08:00:00 (8.00)      $80.00
2      2024-05-02     08:00:00 (8.00)      $80.00
3      2024-05-03     08:00:00 (8.00)      $80.00
4      2024-05-06     08:00:00 (8.00)      $80.00
5      2024-05-07     08:00:00 (8.00)      $80.00
6      2024-05-08     08:00:00 (8.00)      $80.00
7      2024-05-09     08:00:00 (8.00)      $80.00
8      2024-05-10     08:00:00 (8.00)      $80.00
9      2024-05-13     08:00:00 (8.00)      $80.00
10     2024-05-14     08:00:00 (8.00)      $80.00
11     2024-05-15     08:00:00 (8.00)      $80.00
12     2024-05-16     08:00:00 (8.00)      $80.00
13     2024-05-17     08:00:00 (8.00)      $80.00
14     2024-05-20     08:00:00 (8.00)      $80.00
15     2024-05-21     08:00:00 (8.00)      $80.00
16     2024-05-22     08:00:00 (8.00)      $80.00
17     2024-05-23     08:00:00 (8.00)      $80.00
18     2024-05-24     08:00:00 (8.00)      $80.00
Total Hours: 144:00:00 (144.00)

Exchange rate right now for 1 USD to PHP: 58.08
Total Earnings: $1440 (₱83635.20)
Successfully created an invoice: out/invoice-2024-05-27.xlsx
```
