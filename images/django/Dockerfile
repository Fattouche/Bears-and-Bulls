FROM python

RUN mkdir /startup
COPY . /startup
WORKDIR /startup

RUN pip install --upgrade pip
RUN pip install -r /startup/requirements.txt

WORKDIR /day_trader

ENTRYPOINT ["bash", "/startup/entrypoint.sh"]
