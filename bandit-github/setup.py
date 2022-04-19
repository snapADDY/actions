#!/usr/bin/env python
from setuptools import setup, find_packages

setup(
    name='bandit-github',
    version='1.0.0',
    packages=find_packages(include=['bandit_github']),
    install_requires=[
        "bandit",
        "requests"
    ],
    entry_points={'bandit.formatters': ['github = bandit_github.formatter:report']}
)
