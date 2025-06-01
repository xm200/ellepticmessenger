EMess
---

It is decentralized, secure and little bit anonymous messenger.

# Instruction.

  - 1 You must register and then login, after that you are getting a keypair.
    - 1.1 The public key you can share, the private key you keep in secret
  - 2 From server you get the list of users, which is online (only nicknames)
    - 2.1 You choose user, which you want to chat
    - 2.2 You ask server to connect you and that user
    - 2.3 Server asking this user about to chat with you
    - 2.4.1 If the answer is positive, you need to proceed public key exchange.
      - 2.4.1.1 You send user your public key
      - 2.4.1.2 You recieve it`s public key
      - 2.4.1.3 You encrypt message with the user`s public key
      - 2.4.1.4 You send encrypted message via server
      - 2.4.1.5 User decrypt it with it`s private key
      - 2.4.1.6 User making steps 3 and 4, but for you
        
      Great! You are chatting
    - 2.4.2 If the answer is negative.
      - 2.4.2.1 Sorry but you can`t chat
  - 3 If you want to chat with many users at the same time.
    - 3.1 You must have a server.
    - 3.2 The server ask users you have invited about chatting
    - 3.3 The server getting agreed users public keys, and send it`s public key to them
    - 3.4 They send encrypted messages (with server public key) to it, and it is decrypted by it
    - 3.5 Server encrypt message for each user separately and send it
      
    Great, you are chatting with many people!
