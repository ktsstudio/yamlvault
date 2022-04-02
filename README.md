# yamlvault

Utility to transform yaml document containing references to secrets in
HashiCorp Vault by their respective values and generate a new yaml document.

Possible usages include transforming of Helm values.

## Syntax

If you want a value to be replaced it must be of the following format:

```yaml
key: vault:<vault_path>:<vault_key>
```

It also can be an array element:
```yaml
some_array_key:
  - value1
  - vault:<vault_path>:<vault_key>
  - value3
```

where `vault_path` is the path where secret can be fetched and `vault_key` - 
key in the Vault secret 

## Usage



### Decrypt
```text
Usage:
  yamlvault decrypt [flags]

Flags:
  -h, --help            help for decrypt
  -i, --input string    Input file (default "-")
  -o, --output string   Output file (default "-")
```

values.yaml:
```yaml
env:
  POSTGRES_USER: vault:/secret/data/myapp/pg:username
  POSTGRES_PASSWORD: vault:/secret/data/myapp/pg:password
  POSTGRES_DATABASE: vault:/secret/data/myapp/pg:db
  
  MONGO_USER: vault:/mongodb/creds/mymongo:username
  MONGO_PASSWORD: vault:/mongodb/creds/mymongo:password
```

Command:
```bash
% export VAULT_TOKEN="your-vault-token"
$ yamlvault decrypt --input ./vaules.yaml --output ./values.yaml.dec  
```

You can also read directly from stdin and write to stdout (which is actually a default behaviour):
```bash
$ cat values.yaml | yamlvault decrypt > values.yaml.dec
```

values.yaml.dec:
```yaml
env:
  POSTGRES_USER: myapp_user
  POSTGRES_PASSWORD: superpaswword
  POSTGRES_DATABASE: myapp

  MONGO_USER: v-approle-mymongo-we1PN0hdcx3JcLT-16489110
  MONGO_PASSWORD: E-OrMjjgln...
```

As you can see you can also use dynamic secrets (`MONGO_USER` & `MONGO_PASSWORD`)

## Docker

Also available as a docker image:

```bash
$ docker run ktshub/yamlvault 
```