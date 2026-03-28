FROM python:3.11-slim

WORKDIR /workspace

COPY app/requirements.txt /workspace/requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

COPY . /workspace
