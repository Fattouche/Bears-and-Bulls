FROM python

RUN mkdir /service
COPY . /service
WORKDIR /service

RUN pip install --upgrade pip
RUN pip install -r requirements.txt

ENTRYPOINT ["bash", "./entrypoint.sh"]
