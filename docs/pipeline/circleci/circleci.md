Pipeline: CircleCI
==================

**Environment Variables**

* **M8S_API** - The endpoint of the M8s API endpoint eg `m8s.example.com:443`
* **M8S_TOKEN** - The token used for authenticating with the M8s API
* **M8S_GIT_REPO** - A checkout url used for cloning the applications code base eg `https://github.com/previousnext/example`

* **M8S_ENV_FOO** - An additional environment variable which will be injected into the ephermeral environment as "FOO"

**Config**

An example CircleCI environment is available in this folder called `config.yml`

This file demonstrates structure for:

* Running unit tests first
* Bootstrapping the environment
* Executing another step in the new environment