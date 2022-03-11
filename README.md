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

values.yaml:
```yaml
env:
  POSTGRES_USER: vault:/secret/data/myapp/pg:username
  POSTGRES_PASSWORD: vault:/secret/data/myapp/pg:password
  POSTGRES_DATABASE: vault:/secret/data/myapp/pg:db
  
  RABBITMQ_USER: vault:/secret/data/myapp/rabbitmq:username
  RABBITMQ_PASS: vault:/secret/data/myapp/rabbitmq:password
```

Command:
```bash
$ yamlvault decrypt --input ./vaules.yaml --output ./values.yaml.dec  
```

values.yaml.dec:
```yaml
env:
  POSTGRES_USER: myapp_user
  POSTGRES_PASSWORD: superpaswword
  POSTGRES_DATABASE: myapp
  
  RABBITMQ_USER: myapp_user
  RABBITMQ_PASS: rmq_super_password
```
