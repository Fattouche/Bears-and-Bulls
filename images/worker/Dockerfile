FROM python

RUN mkdir /service
COPY . /service
WORKDIR /service

COPY ./requirements.txt .
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

COPY . /service

ENTRYPOINT ["bash", "./images/worker/entrypoint.sh"]
