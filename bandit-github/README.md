# Bandit github

Posts [bandit](https://github.com/PyCQA/bandit) security issues as a comment in a pull request.

## Test the workflow locally

To test the workflow locally, one could use the following workflow.


### 1. Create an environment and install module

Navigate to `bandit-github`, create a virtual environment and install the package.

```sh
>>> python3.9 -m venv .env
>>> source .env/bin/activate
>>> pip install . --force-reinstall
```

### 2. Create an example directory

Within `bandit-github`, create a directory `examples` and a test file `example.py`.

```sh
>>> mkdir examples
>>> cd examples
>>> echo import lxml > example.py
```

### 3. Run bandit

Run bandit by using the `entrypoint.sh` script within the `bandit-github` directory. The `examples/example.py` script will return an issue.

```sh
>>> source ./entrypoint.sh . UNDEFINED UNDEFINED ./.env 0 DEFAULT DEFAULT
[main]  INFO    profile include tests: None
[main]  INFO    profile exclude tests: None
[main]  INFO    cli include tests: None
[main]  INFO    cli exclude tests: None
[main]  INFO    running on Python 3.9.5
1. Issue: 'Using lxml to parse untrusted XML data is known to be vulnerable to XML attacks. Replace lxml with the equivalent defusedxml package.' from B410:lxml: CWE: CWE-20 (https://cwe.mitre.org/data/definitions/20.html), Severity: LOW Confidence: HIGH at ./examples/example.py:1
```
