# txtme

```sh
run-a-long-script && txtme
```

## Setup

1. Create a Twilio dev project -> [here](https://www.twilio.com/try-twilio)
    
2. Add a `.txtme.toml` like this

```toml
Name = "<your name>"
PhoneTo = "<your phone>"
PhoneFrom = "<twilio project phone>"
SID = "<twilio sid>"
Token = "<twilio token>"
```

If you do not add a .txtme.toml file, the cli will guide you through creating one on first run.
