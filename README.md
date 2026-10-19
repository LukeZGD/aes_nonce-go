# aes_nonce-go

Go implementations of [Cryptiiiic/aes_nonce](https://github.com/Cryptiiiic/aes_nonce), including `aes_nonce` and `aes_cryptex_nonce`. Written in Go 1.14-compatible code with no external dependencies.

## Building

```bash
go build -o aes_nonce aes_nonce.go
go build -o aes_cryptex_nonce aes_cryptex_nonce.go
```

## Usage

### 1. Get your device keys

Prerequisites:

- Install [x8A4](https://github.com/Cryptiiiic/x8A4) on your jailbroken iOS device.
- You can install through Sileo or Zebra using this repo: https://lukezgd.github.io/x8a4
- Installing OpenSSH and connecting to your device over SSH is recommended to make copying the key output easier.

Run the appropriate command on your device.

**For the ApNonce:**

```bash
sudo x8A4 -k 0x8A3
```

Save the `0x8A3` key.

**For the Cryptex nonce:**

```bash
sudo x8A4 -k 0x8A4
```

Save the `0x8A4` key.

### 2. Generate the entangled nonces

Run the appropriate executable with your key and preferred generator. **Omit the `0x` prefix from the key and generator arguments.**

**ApNonce:**

```bash
./aes_nonce <0x8A3-key> <generator>
```

Example:

```bash
./aes_nonce 4F6A4BA34839AC1D5764CC32871B19A4 1111111111111111
```

**Cryptex nonce:**

```bash
./aes_cryptex_nonce <0x8A4-key> <cryptex-generator>
```

Example:

```bash
./aes_cryptex_nonce 2CDD86786C588EC35CCF68273FD44685 11111111111111111111111111111111
```

Keep the encrypted generator/seed and entangled nonce outputs for later use.

## Saving SHSH Blobs

**Important:** The following instructions cover saving blobs for iOS versions that Apple is currently signing. This process is separate from saving onboard blobs.

- **Currently signed versions:** Use `tsschecker` with the appropriate generator and entangled nonce values, as described below.
- **Onboard blobs:** This is a separate process that does not require `aes_nonce` or `aes_cryptex_nonce`. However, [x8A4](https://github.com/Cryptiiiic/x8A4) is still required. See here: https://github.com/LukeZGD/Legacy-iOS-Kit/wiki/Saving-onboard-SHSH-blobs-of-current-iOS-version

For additional information on saving Cryptex1 blobs, see this guide: https://ichitaso.com/ios/shsh/how-to-save-shsh-for-cryptex1/

### Saving blobs for currently signed versions

Use [1Conan/tsschecker](https://github.com/1Conan/tsschecker) to save SHSH blobs.

Connect your device in normal mode and run the following command as a reference. Replace `27.0.1` with the desired iOS version, and replace the generator and Cryptex generator values with your chosen values.

```bash
./tsschecker \
    -d "$(ideviceinfo -k ProductType)" \
    -i 27.0.1 \
    -s \
    -e "$(ideviceinfo -k UniqueChipID)" \
    -g 0x1111111111111111 \
    -x 0x11111111111111111111111111111111 \
    -8 <entangled-nonce> \
    -t <entangled-cryptex-nonce> \
    -c "$(ideviceinfo -k BasebandSerialNumber)"
```

Replace `<entangled-nonce>` and `<entangled-cryptex-nonce>` with the outputs from `aes_nonce` and `aes_cryptex_nonce` respectively.

**Notes:**

- `-d` specifies the device product type.
- `-i` specifies the iOS version.
- `-e` specifies the device ECID.
- `-g` specifies the generator.
- `-x` specifies the Cryptex generator.
- `-8` specifies the entangled nonce.
- `-t` specifies the entangled Cryptex nonce.
- Omit `-c` for devices without a baseband, such as Wi-Fi-only iPads.
- Make sure the generator and entangled nonce values correspond to the values you intend to use.
