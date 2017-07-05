#!/usr/bin/env python
# -*- coding: utf-8 -*-

from flask import Flask, render_template, request

app = Flask(__name__)

@app.route('/')
def root():
    return render_template("hello.html")

