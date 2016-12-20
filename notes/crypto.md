Cryptanalytical capabilities. Is the math there?
1. Elliptic curves - smaller than finite field, less well understood, can a room full of math people invent new math to break them? classes of curves are more vul, nsa has tried to affect curve selection.
2. Factoring discreate logrithms
3. rc4 - barely secure, used in a lot of internet protocols.

Most crypto causes nsa trouble, at least at scale, so they get around it: default weak keys, sabotage standards, insert back doors, exfiltrating keys. VPN systems are hosed? mostly relies on unencrypted streams of data, or cooperation with companies. Bulk collection is too easy.

We're now in a world where there are pr benefits to fighting, reputation and hence sales are at stake. Previously the cost of cooperating with NSA was low.

TextSecure protocol from open whisper systems

Threat model
prevent mitm in a forward secret manner
Bob publishes a public key, Alice pulls it
How to encrypt and authenticate message?
How does Alice know the public key is really Bob's public key?

Simple public/private key doesn't preserve forward secrecy?
1. prekeys:
    Bob generates 100 prekeys, alice picks one, Bob deletes key upon receiving the message
    Bob might run out of prekeys on DoS
    If directory is compromised, Bob can't decrypt the message

2. Rotated prekeys:
    Periodically sign new key and replace old key
    not forward secure for the given time period

3. Both
    both Bob and Alice have a long term identity key
    Alice sends hers in first message
    Create a shared session key with forward security



