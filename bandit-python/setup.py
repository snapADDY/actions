#!/usr/bin/env python
from setuptools import find_packages, setup

setup(
    name="bandit-python",
    version="1.0.0",
    packages=find_packages(include=["bandit_python"]),
    install_requires=["bandit", "requests"],
    entry_points={"bandit.formatters": ["github = bandit_python.formatter:report"]},
)
