#!/usr/bin/env bash
./gohack Add.asm | diff Add.hack -
./gohack MaxL.asm | diff MaxL.hack -
./gohack Max.asm | diff Max.hack -
./gohack Rect.asm | diff Rect.hack -
./gohack Pong.asm | diff Pong.hack -
