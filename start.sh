#!/bin/bash

./room &
cd web/chat && yarn serve --port 8080 &