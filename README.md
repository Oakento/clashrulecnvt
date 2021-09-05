# clashrulecnvt
This tool converts `rule-provider` field and `RULE-SET` announcement that [clash premium](https://github.com/Dreamacro/clash/releases/tag/premium) proprietorially supports into simple rule entries.


## Configure
Edit the default configuration `$HOME/.config/clashrulecnvt/config.yaml`.

## Usage
Use the default configuration file:
```bash
$ clashrulecnvt
```
Specify a configuration file:
```bash
$ clashrulecnvt -c /path/to/file
```
Arguments specified in commandline will be replace the value in the configuration.
```sh
$ clashrulecnvt -i /path/to/input -o /path/to/output -p http://127.0.0.1:7890
```
