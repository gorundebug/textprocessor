#!/bin/bash


if ! docker build -t textprocessor .; then
    exit 1
fi



exit 0