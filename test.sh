#!/usr/bin/env bash
./gohack Pong.asm | diff Pong.hack -
./gohack Add.asm | diff Add.hack -
./gohack Max.asm | diff Max.hack -
./gohack MaxL.asm | diff MaxL.hack -
./gohack Rect.asm | diff Rect.hack -
