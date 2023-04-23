#!/bin/bash

set -e

./room &
cd web && yarn serve --port 8080
