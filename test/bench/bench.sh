#!/usr/bin/env sh

ab -H Host:bench.test $@ http://reverse-proxy/ping
ab -H Host:bench.test $@ http://nginx/ping
ab -H Host:bench.test $@ http://apache/ping

