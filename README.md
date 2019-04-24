# Shadow APM
Shadow APM is the a request simulator for a microservice architecture. It stores all of the request signature in a database(MongoDB), which can be simulated later. Elastic APM agent is need to be integrated and hence configured to capture the request body.
This project is compatible with all of the [APM Agents by Elastic](https://www.elastic.co/guide/en/apm/agent/index.html).
Shadow APM can run individually or alongside [APM Server by Elastic](https://github.com/elastic/apm-server) by shadowing the traffic using tools like  [GoReplay](https://github.com/buger/goreplay) or [GoDuplicator](https://github.com/mkevac/goduplicator).

## WORK ON PROGESS
There are lots of improvement that can be done. Help will be really appreciated. Feel free to send a PR.

You can also reach me at `coolboi567@gmail.com`.
