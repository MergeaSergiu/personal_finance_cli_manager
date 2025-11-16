# Personal_Finance_Cli_Manager

The **Personal Finance CLI Manager** is a command-line
application developed in Go.  
Its purpose is to facilitate the management, analysis, and categorization of
personal financial transactions through an efficient, text-based interface, while utilizing an SQLite database to ensure persistent, reliable, and structured storage of all financial data.

##  Tech Stack

- Go: 1.25.4
- Database: SQLite
- CLI Framework: Bubbletea https://github.com/charmbracelet/bubbletea

##  Features

- A text-based terminal user interface implemented using Bubbletea. 

- The user can import financial transactions from CSV or OFX files.

- The user can manually add income or expense transactions.

- The system can categorize transactions automatically based on user-defined
  rules (e.g., using regular expressions).

- The system allows budget tracking with alerts

- The user can generate various reports containing transaction information

- The user can search transactions and filter data using specific criteria.

- System ensures data consistency and integrity
