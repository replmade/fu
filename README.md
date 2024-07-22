# fu

A CLI tool to authenticate a Firebase user with email/password and generate tokens.

Tokens that can be generated:
- id token
- session cookie
- custom token

Depends on [firebase-spells-go](https://github.com/replmade/firebase-spells-go) for all Firebase operations

## Commands

**init**: Initializes a Firebase application
```bash
./fu --app-name <my-app> --api-key <firebase api key> --sa-key-path <service account key file path>
```

Saves the Firebase API key and service account key file path to $HOME/.fu/config.toml
```toml
[my-app]
api_key = <firebase api key>
sa_key_path = <service account key file path>
```

**load**: Loads a Firebase application
```bash
./fu load --app-name <my-app>
```

**signin**: Attemps to sign a user into the currently loaded Firebase app with a user's email and password. If successful, an ID token is returned from Firebase
```bash
./fu signin --email <user@email> --password <user-password>
```

**id-token**: After the **signin** command, use this command to print the id token. This token is also stored in the config file under the currently active Firebase app section
```bash
./fu id-token
```

**session**: After the **signin** command, this command prints a session token for the user. Like the id token, this is stored in the config.toml. There's an optional argument to specify expiry in seconds.
```bash
./fu session --expires-in <integer-value>
```

