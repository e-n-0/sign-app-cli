# Sign-App-CLI

Sign-App-CLI is a command-line tool for signing iOS and macOS apps. This tool is built to make the code-signing process quick and easy, saving developers precious time that can be spent on other important tasks.
This tool can be integrated into CI/CD pipelines to automate the code-signing process.

This tool is written in Go but can be used only on macOS or iOS (jailbroken) devices.

## Features

- Sign iOS and MacOS apps from the command line (.ipa files and .app folders)
- Bundle apps into ipa files that are ready to be installed on iPhones or Silicon Macs
- List codesigning certificates installed on your machine
- List provisioning profiles installed on your machine

## How to get it

### From Releases

You can download the latest version of the tool from the [releases page](https://github.com/e-n-0/sign-app-cli/releases).

### Compile from source

You can also compile the tool from source. You will need to have Go installed on your machine.

```bash
git clone git@github.com:e-n-0/sign-app-cli.git
cd sign-app-cli
go build
```

## Usage

```bash
sign-app-cli -h
```
```
Sign your iOS/Macos app from command line

Usage:
  sign-app-cli [command]

Available Commands:
  help                     Help about any command
  listCodesigningCerts     List all codesigning certificates available in your keychain
  listProvisioningProfiles List all provisioning profiles available in your keychain
  sign                     Sign the provided file

Flags:
  -h, --help   help for sign-app-cli

Use "sign-app-cli [command] --help" for more information about a command.
```

### List codesigning certificates

```bash
sign-app-cli listCodesigningCerts
```
```
Found 3 codesigning certificates:
   Apple Development: Fake Person (XXXXXXXXXX)
   Apple Development: Imaginary Name (YYYYYYYYYY)
   Apple Distribution: Fake Company (ZZZZZZZZZZ)
```

### List provisioning profiles

```bash
sign-app-cli listProvisioningProfiles
```
```
Found 2 provisioning profiles:
  MyMobileProvision (XXXXXXXXXX)
  Test (YYYYYYYYYY)
```

### Sign an app

```bash
sign-app-cli sign -h
```
```
Sign the provided file with the provided provisioning profile and codesigning certificate.
The file can be an .ipa or a .app.

Usage:
  sign-app-cli sign [flags]

Flags:
  -c, --certificate string    The name of the codesigning certificate to use installed on the machine
  -e, --entitlements string   The path of the entitlements file to use
  -h, --help                  help for sign
  -i, --input string          The path of the file to sign
  -o, --output string         The path of the signed file
  -p, --profile string        The name of the provisioning profile to use installed on the machine
  -P, --profilePath string    The path of the provisioning profile to use
```

### Example

I want to sign the app located at `/Users/fakeperson/Desktop/MyApp.ipa` with the provisioning profile `MyMobileProvision (XXXXXXXXXX)` and the certificate `Apple Development: Fake Person (XXXXXXXXXX)`.

```bash
sign-app-cli sign -i /Users/fakeperson/Desktop/MyApp.ipa -p "MyMobileProvision (XXXXXXXXXX)" -c "Apple Development: Fake Person (XXXXXXXXXX)" -o /Users/fakeperson/Desktop/MyApp-signed.ipa
```

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.