#!/bin/bash

./room &
cd web && yarn serve --port 8080 &