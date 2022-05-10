FROM python:3.9-alpine

WORKDIR /app

COPY bandit_github bandit_github
COPY setup.py setup.py
RUN pip install .

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
