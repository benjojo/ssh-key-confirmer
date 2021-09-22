ssh-key-confirmer
===

This confirms if a SSH public key is listed as a authorized_key on a system

## Usage

```
$ ssh-key-confirmer -i ./id_rsa.pub git@github.com
Key is present on user+server
$ ssh-key-confirmer -i ./id_ed25519.pub ben@localhost
Key not found on user+server
```

## How does this work

If you ssh with verbose mode enabled, you can see that the way ssh works is that you
offer a key to a server, and it will respond back if it accepts or not (and then you
provide a signed copy of a shared session secret as proof you own the key)

```
$ ssh -v -i id_ed25519 localhost
OpenSSH_8.2p1 Ubuntu-4ubuntu0.3, OpenSSL 1.1.1f  31 Mar 2020
debug1: Reading configuration data /home/ben/.ssh/config
...
debug1: Connecting to localhost [127.0.0.1] port 22.
debug1: Connection established.
...
debug1: Authenticating to localhost:22 as 'ben'
debug1: SSH2_MSG_KEXINIT sent
debug1: SSH2_MSG_KEXINIT received
...
debug1: Authentications that can continue: publickey,password
debug1: Next authentication method: publickey
debug1: Offering public key: cardno:000605032939 RSA SHA256:CtGA1RT0bAOd06HHIRQTB9reCED/SrD2MGfS8MUEd6Q agent
debug1: Authentications that can continue: publickey,password
>>>>>> debug1: Offering public key: id_ed25519 ED25519 SHA256:Wml16ewzvx7SBMLmFkvVxZBiwN5lAcFm6nuLJF2rKYY explicit <<<<<<
>>>>>> debug1: Server accepts key: id_ed25519 ED25519 SHA256:Wml16ewzvx7SBMLmFkvVxZBiwN5lAcFm6nuLJF2rKYY explicit  <<<<<<
debug1: Authentication succeeded (publickey).
Authenticated to localhost ([127.0.0.1]:22).
...
Welcome to Ubuntu 20.04.3 LTS (GNU/Linux 5.10.26-2fast2benjojo x86_64)
```

This tool simply offers the key, and if it's accepted it will confirm that key exists, a server could be configured to
accept every possible key. There is an attempt to detect when this is happening.

## How do I know this is being done to me?

You will get `[preauth]` disconnections, though a lot of other things can cause that too so... /shrug

```
Connection closed by authenticating user ben 127.0.0.1 port 59566 [preauth]
```

## That's a janky name

Well, `ssh-keyscan` was taken, and I'm pretty sure i've seen a `ssh-check-key` somewhere.
