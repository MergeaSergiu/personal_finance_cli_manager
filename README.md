# Personal_Finance_Cli_Manager

The **Personal Finance CLI Manager** is a command-line
application developed in Go.  
Its purpose is to facilitate the management, analysis, and categorization of
personal financial transactions through an efficient, text-based interface, while utilizing an SQLite database to ensure persistent, reliable, and structured storage of all financial data.

Provides category management, transactions, CSV import, monthly expense charts, budget overview and an email notification system using RabbitMQ and SMTP (Maildev for local testing). The TUI is powered by Bubble Tea.

## Tech Stack

- Go: 1.25.4
- Database: SQLite
- CLI / TUI: Bubble Tea (https://github.com/charmbracelet/bubbletea)
- Queue: RabbitMQ
- Local SMTP (development): Maildev

## Features

- Import transactions from CSV
- Manually add income and expense transactions
- Manually add expense category
- Automatic categorization using user-defined rules (e.g., regex)
- Budget tracking with alerts
- The system generates charts for budget spendings overview
- Generates reports for monthly spendings
- The user can search & filter transactions
- Email notifications queued via RabbitMQ and delivered via SMTP


## How to run the project

1. Start required services from `docker-compose.yml` first. This brings up RabbitMQ and Maildev required by the email/queue subsystem.

   - Docker Compose V2 (recommended):
     ```powershell
     docker compose -f docker-compose.yml up -d
Confirm containers are running before starting the application.

2. Build or run the application:

   - Build:
     ```powershell
     go build ./...
     .\peronal_finance_cli_manager.exe
     ```

   - Run (development):
     ```powershell
     go run ./cmd/main.go
     ```

## Install & Build

Clone the repository and build:

```powershell
git clone https://github.com/MergeaSergiu/peronal_finance_cli_manager.git
cd peronal_finance_cli_manager
go build ./...
 ```

## Personal Finance Manager UI

<img width="1847" height="385" alt="ss1" src="https://github.com/user-attachments/assets/8b88f77f-cbc3-42cc-b516-c0373063bb3c" />

<img width="1468" height="587" alt="ss2" src="https://github.com/user-attachments/assets/c4f5071b-7bcc-4e16-ba7f-467268420c6f" />

<img width="1827" height="350" alt="ss3" src="https://github.com/user-attachments/assets/75952271-2eaf-48e7-8a74-b1af89b9491c" />

<img width="1917" height="792" alt="ss4" src="https://github.com/user-attachments/assets/9fec48fd-831a-4972-8488-de9a89025521" />

<img width="1871" height="551" alt="ss5" src="https://github.com/user-attachments/assets/8a8a48a0-64c3-4088-88bb-e01af87e3a59" />

<img width="1892" height="588" alt="ss6" src="https://github.com/user-attachments/assets/20f9b592-1690-4918-bc34-5075bb3468e8" />

<img width="1836" height="432" alt="ss7" src="https://github.com/user-attachments/assets/cf1d225f-b0b9-44e0-9a52-dcfd5a0b23b1" />

<img width="1892" height="380" alt="ss8" src="https://github.com/user-attachments/assets/1b61269a-49aa-4bf0-9535-55badc65114c" />

<img width="1858" height="395" alt="ss9" src="https://github.com/user-attachments/assets/963ffadb-cc0b-44f8-8c9e-b1dddbd62080" />

<img width="1862" height="326" alt="ss10" src="https://github.com/user-attachments/assets/14b5a104-dc96-4679-bf61-176c3306ff88" />

<img width="1911" height="501" alt="ss11" src="https://github.com/user-attachments/assets/2af897f9-ce14-4580-9d69-82713f2d2f67" />
















    


  

     


