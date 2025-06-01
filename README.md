# DNSwitch

DNSwitch is a cross-platform command-line tool for easily managing and switching between different DNS servers. It supports Windows, Linux, and macOS, allowing users to quickly change their DNS settings, add custom DNS servers, and clear DNS configurations.

## Features

- Switch between predefined DNS servers (403, Shecan, Begzar, Electrotm)
- Add and remove custom DNS servers
- Clear all DNS settings
- Automatically detect the active network interface
- Cross-platform support (Windows, Linux, macOS)
- Persistent storage of custom DNS servers

## Requirements

- Go 1.23.4
- Administrator/root privileges (for changing DNS settings)

## Installation

You can install DNSwitch using one of the following methods:

### Method 1: Using go install (Recommended)

```
go install github.com/naghinezhad/DNSwitch@latest
```

This will install the latest version of DNSwitch directly to your `$GOPATH/bin` or `$GOBIN` directory.

### Method 2: Build from source

1. Clone the repository:

   ```
   git clone https://github.com/naghinezhad/DNSwitch.git
   ```

2. Navigate to the project directory:

   ```
   cd DNSwitch
   ```

3. Build the program:
   ```
   go build
   ```

## Usage

Run the program with administrator/root privileges:

- On Windows:
  Right-click on the Command Prompt or PowerShell and select "Run as administrator", then run:

  ```
  DNSwitch
  ```

  If you built from source, navigate to the program directory and run:
  ```
  .\DNSwitch.exe
  ```

- On Linux/macOS:
  ```
  sudo DNSwitch
  ```

  If you built from source, navigate to the program directory and run:
  ```
  sudo ./DNSwitch
  ```

## Main Menu Options

1. Switch to a predefined or custom DNS server
2. Add a custom DNS server
3. Remove a custom DNS server
4. Clear all DNS settings
5. Exit the program

## Adding a Custom DNS Server

1. Select "Add custom DNS" from the main menu
2. Enter a name for the DNS server
3. Enter the primary IP address
4. Enter the secondary IP address (optional)

## Removing a Custom DNS Server

1. Select "Remove custom DNS" from the main menu
2. Choose the custom DNS server you want to remove from the list
3. Confirm the removal

## Platform-Specific Notes

### Windows

- Uses PowerShell commands to manage DNS settings
- Requires running the program as an administrator

### Linux

- Modifies the `/etc/resolv.conf` file to change DNS settings
- Requires sudo privileges
- Note: Some Linux distributions may use different methods for DNS configuration (e.g., NetworkManager). This program might not work correctly in such cases.

### macOS

- Uses the `networksetup` command to manage DNS settings
- Requires sudo privileges

## Important Notes

- Always run the program with administrator/root privileges to ensure it can modify system DNS settings.
- Custom DNS servers are stored in a file named `custom_dns.json` in the same directory as the program. This file is created automatically when you add a custom DNS server.
- The program detects the active network interface automatically. If you have multiple active interfaces, it may not always select the desired one.
- Changes to DNS settings may not take effect immediately on all systems. You may need to flush your DNS cache or restart your network connection for changes to apply.

## Limitations

- The program does not support IPv6 DNS configurations.
- On Linux, the program assumes that `/etc/resolv.conf` is used for DNS configuration, which may not be true for all distributions or configurations.
- The program does not provide options for advanced DNS configurations (e.g., separate DNS servers for different domains).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request to the [DNSwitch repository](https://github.com/naghinezhad/DNSwitch).

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgements

DNSwitch was created to simplify the process of changing DNS settings across different platforms. It aims to provide a user-friendly interface for managing DNS configurations, especially for users in regions where changing DNS servers is common practice.

## Support

If you encounter any issues or have questions, please open an issue on the [GitHub repository](https://github.com/naghinezhad/DNSwitch/issues).

Thank you for using DNSwitch!